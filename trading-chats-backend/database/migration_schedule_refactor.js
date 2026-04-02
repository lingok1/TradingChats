const dbName = db.getName()
const systemConfigId = 'global_config'

function log(message) {
  print(`[migration_schedule_refactor] ${message}`)
}

function getDistinctPromptValues(batchId) {
  if (!batchId) {
    return []
  }
  return db.ai_responses.distinct('prompt', { batch_id: batchId }).filter(Boolean)
}

function collectScheduleConfigParamCombos() {
  return db.schedule_configs.aggregate([
    {
      $project: {
        _id: 0,
        param1: { $ifNull: ['$param1', ''] },
        param2: { $ifNull: ['$param2', ''] },
      },
    },
    {
      $group: {
        _id: { param1: '$param1', param2: '$param2' },
        count: { $sum: 1 },
      },
    },
  ]).toArray()
}

function migrateUp() {
  log(`start up migration on ${dbName}`)

  const backupCollection = `schedule_logs_backup_${new Date().toISOString().replace(/[.:TZ-]/g, '')}`
  db.schedule_logs.aggregate([{ $match: {} }, { $out: backupCollection }])
  log(`schedule_logs backup created: ${backupCollection}`)

  const combos = collectScheduleConfigParamCombos()
  if (combos.length > 1) {
    log(`conflicting schedule param combinations found: ${tojson(combos)}`)
    throw new Error('multiple param1/param2 combinations exist in schedule_configs; aborting destructive migration')
  }

  const runtimeConfig = combos.length === 1 ? combos[0]._id : { param1: '', param2: '' }
  db.system_configs.updateOne(
    { _id: systemConfigId },
    {
      $setOnInsert: {
        _id: systemConfigId,
        system_title: 'Trading Chats',
        system_logo: '',
        parameters: {},
      },
      $set: {
        param1: runtimeConfig.param1 || '',
        param2: runtimeConfig.param2 || '',
        updated_at: new Date(),
      },
    },
    { upsert: true }
  )
  log(`system_configs runtime params set to param1=${runtimeConfig.param1 || ''}, param2=${runtimeConfig.param2 || ''}`)

  let migratedLogs = 0
  let missingBatch = 0
  let inconsistentPrompt = 0

  db.schedule_logs.find({}).forEach((doc) => {
    const prompts = getDistinctPromptValues(doc.batch_id)
    const patch = {
      trigger_type: doc.trigger_type || 'auto',
    }

    if (!doc.prompt) {
      if (prompts.length === 1) {
        patch.prompt = prompts[0]
        migratedLogs += 1
      } else if (prompts.length === 0) {
        missingBatch += 1
      } else {
        inconsistentPrompt += 1
      }
    }

    db.schedule_logs.updateOne({ _id: doc._id }, { $set: patch })
  })

  db.schedule_configs.updateMany({}, { $unset: { param1: '', param2: '' } })
  log('removed param1/param2 from schedule_configs documents')

  try {
    db.schedule_configs.dropIndex('param1_1')
    log('dropped index param1_1')
  } catch (e) {
    log(`skip drop index param1_1: ${e.message}`)
  }

  try {
    db.schedule_configs.dropIndex('param2_1')
    log('dropped index param2_1')
  } catch (e) {
    log(`skip drop index param2_1: ${e.message}`)
  }

  log(`schedule_logs prompt migrated=${migratedLogs}, missingBatch=${missingBatch}, inconsistentPrompt=${inconsistentPrompt}`)
  log('up migration complete')
}

function migrateDown() {
  log(`start down migration on ${dbName}`)

  const systemConfig = db.system_configs.findOne({ _id: systemConfigId }) || {}
  const param1 = systemConfig.param1 || ''
  const param2 = systemConfig.param2 || ''

  db.schedule_configs.updateMany(
    {},
    {
      $set: {
        param1,
        param2,
      },
    }
  )
  log('restored param1/param2 onto schedule_configs from system_configs')

  db.schedule_logs.updateMany(
    {},
    {
      $unset: {
        prompt: '',
        trigger_type: '',
      },
    }
  )
  log('removed prompt/trigger_type from schedule_logs')

  db.system_configs.updateOne(
    { _id: systemConfigId },
    {
      $unset: {
        param1: '',
        param2: '',
      },
      $set: {
        updated_at: new Date(),
      },
    }
  )
  log('removed runtime params from system_configs')

  log('down migration complete')
}

function validate() {
  log(`start validation on ${dbName}`)

  let ok = true
  let checked = 0
  let mismatched = 0
  let missingTriggerType = 0

  db.schedule_logs.find({}).forEach((doc) => {
    checked += 1
    const prompts = getDistinctPromptValues(doc.batch_id)
    if (prompts.length === 1 && (doc.prompt || '') !== prompts[0]) {
      ok = false
      mismatched += 1
      log(`prompt mismatch for schedule_log=${doc._id}, batch_id=${doc.batch_id}`)
    }
    if (!doc.trigger_type || ['manual', 'auto'].indexOf(doc.trigger_type) === -1) {
      ok = false
      missingTriggerType += 1
      log(`invalid trigger_type for schedule_log=${doc._id}`)
    }
  })

  const remainingParamDocs = db.schedule_configs.countDocuments({
    $or: [{ param1: { $exists: true } }, { param2: { $exists: true } }],
  })
  if (remainingParamDocs > 0) {
    ok = false
    log(`schedule_configs still contains param1/param2 in ${remainingParamDocs} docs`)
  }

  const runtimeConfig = db.system_configs.findOne({ _id: systemConfigId }) || {}
  if (typeof (runtimeConfig.param1 || '') !== 'string' || typeof (runtimeConfig.param2 || '') !== 'string') {
    ok = false
    log('system_configs param1/param2 are not strings')
  }

  log(`validation summary checked=${checked}, mismatched=${mismatched}, invalidTriggerType=${missingTriggerType}`)
  if (!ok) {
    throw new Error('validation failed')
  }
  log('validation complete')
}

const action = typeof MIGRATION_ACTION === 'string' ? MIGRATION_ACTION : 'up'
if (action === 'up') {
  migrateUp()
} else if (action === 'down') {
  migrateDown()
} else if (action === 'validate') {
  validate()
} else {
  throw new Error(`unknown MIGRATION_ACTION: ${action}`)
}
