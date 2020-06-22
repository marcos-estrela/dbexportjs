const connection = require('./connection')

const database = process.env.DATABASE

const FUNCTION = 'FUNCTION'
const PROCEDURE = 'PROCEDURE'

const executeQuery = async (query, params = []) => {
  return new Promise(function(resolve, reject){
      connection.query(query, params, function (error, results, fields){
          if(error){
              let joinedParams = params.map(p => p === null ? 'NULL' : p).join(',')
              let stringToReplace = params.map( p => '?').join(',')
              query = query.replace(/ ,/g, ',').replace(/, /g, ',')
              let completeQuery = query.replace(stringToReplace, joinedParams)
              reject(new Error(`${error} \n\n Executed Query: ${completeQuery}`));
          }
          resolve(results);
      });
  })
}

const end = async () => {
  connection.end()
}

const getParameters = async (objectName, objectType) => {
  const sql = `SELECT par.SPECIFIC_NAME, par.PARAMETER_MODE, par.PARAMETER_NAME, par.DTD_IDENTIFIER, par.CHARACTER_SET_NAME, sch.DEFAULT_CHARACTER_SET_NAME
    FROM INFORMATION_SCHEMA.PARAMETERS par
    JOIN INFORMATION_SCHEMA.SCHEMATA sch ON par.SPECIFIC_SCHEMA = sch.SCHEMA_NAME
    WHERE
    par.SPECIFIC_SCHEMA = ?
    AND par.SPECIFIC_NAME = ?
    AND par.ROUTINE_TYPE = ?
    ORDER BY par.SPECIFIC_NAME, par.ORDINAL_POSITION ASC`
  try {
    const results = await executeQuery(sql, [database, objectName, objectType])
    return results
  } catch (err) {
    console.log(err)
    return []
  }
}

const getTables = async (tableName) => {
  let whereTableName = '';
  if(tableName) {
    whereTableName = `AND TABLE_NAME = ?`
  }
  const sql = `SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = ? AND TABLE_SCHEMA = ? ${whereTableName}`

  try {
    let tables = []
    let params = ['BASE TABLE', database]
    if(tableName){
      params.push(tableName)
    }
    const results = await executeQuery(sql, params)
    tables = await results.map(async result => {
      const name = result['TABLE_NAME']
      const tableResult = await executeQuery(`SHOW CREATE TABLE ${name}`)
      const content = tableResult[0]['Create Table'].replace(/ AUTO_INCREMENT=[0-9]*/gi, '')

      return { name: name, content: content }
    })

    return tables
  } catch (error) {
    console.log(error)
    return []
  }
}

const getFunctions = async (functionName) => {
  let whereFunctionName = ''
  if(functionName){
    whereFunctionName = 'AND ROUTINE_NAME = ?'
  }
  const sql = `SELECT ROUTINE_NAME, ROUTINE_TYPE, ROUTINE_DEFINITION, ROUTINE_COMMENT, DTD_IDENTIFIER
    FROM INFORMATION_SCHEMA.ROUTINES
    WHERE
    ROUTINE_SCHEMA = ?
    AND ROUTINE_TYPE = ?
    ${whereFunctionName}`
  try {
    let functions = []
    let params = [database, FUNCTION]
    if(functionName) {
      params.push(functionName)
    }
    const results = await executeQuery(sql, params)
    functions = await results.map(async result => {
      const name = result['ROUTINE_NAME']
      const charSetName = result['CHARACTER_SET_NAME']
      const charSet = charSetName ? ` CHARSET ${charSetName}` : ''
      let parametersList = []
      let parameters = await getParameters(name, FUNCTION)
      parametersList = await parameters.map(params => makeParametersForFunctions(params))
      const returnType = parametersList.shift().replace('null ', '')
      const content = result['ROUTINE_DEFINITION']
      let functionContent = `CREATE FUNCTION ${name}(${parametersList.join(', ')}) RETURNS ${returnType}${charSet}
${content}`

      return { name: name, content: functionContent }
    })

    return functions
  } catch (error) {
    console.log(error)
    return []
  }
}

const makeParametersForFunctions = (params) => makeParameters(params, FUNCTION)

const getProcedures = async (procedureName) => {
  let whereProcedureName = ''
  if(procedureName){
    whereProcedureName = 'AND ROUTINE_NAME = ?'
  }
  const sql = `SELECT ROUTINE_NAME, ROUTINE_TYPE, ROUTINE_DEFINITION, ROUTINE_COMMENT
    FROM INFORMATION_SCHEMA.ROUTINES
    WHERE
    ROUTINE_SCHEMA = ?
    AND ROUTINE_TYPE = ?
    ${whereProcedureName}`
  try {
    let procedures = []
    let params = [database, PROCEDURE]
    if(procedureName){
      params.push(procedureName)
    }
    const results = await executeQuery(sql, params)
    procedures = await results.map(async result => {
      const name = result['ROUTINE_NAME']
      const comment = result['ROUTINE_COMMENT']
      let parametersList = []
      let parameters = await getParameters(name, 'PROCEDURE')
      parametersList = await parameters.map(params => makeParametersForProcedures(params))
      const content = result['ROUTINE_DEFINITION']
      let procedure = `CREATE PROCEDURE ${name}(${parametersList.join(', ')})
    COMMENT '${comment}'
${content}`

      return { name: name, content: procedure }
    })

    return procedures
  } catch (error) {
    console.log(error)
    return []
  }
}

const makeParametersForProcedures = (params) => makeParameters(params, PROCEDURE)

const makeParameters = (params, routineType) => {
  let charset = '';
  let parameter = '';

  if (routineType === PROCEDURE) {
    parameter = `${params['PARAMETER_MODE']}`
  }

  parameter = `${parameter} ${params['PARAMETER_NAME']} ${params['DTD_IDENTIFIER']}`

  if (params['CHARACTER_SET_NAME'] && (params['CHARACTER_SET_NAME'] !== params['DEFAULT_CHARACTER_SET_NAME'])) {
    charset = `CHARSET ${params['CHARACTER_SET_NAME']}`
  }

  return `${parameter} ${charset}`
}

const getTriggers = async (triggerName) => {
  let whereTriggerName = ''
  if(triggerName){
    whereTriggerName = 'AND TRIGGER_NAME = ?'
  }
  const sql = `SELECT TRIGGER_NAME, ACTION_STATEMENT AS CONTENT, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORIENTATION
    FROM INFORMATION_SCHEMA.TRIGGERS
    WHERE
    TRIGGER_SCHEMA = ?
    ${whereTriggerName}`

  try {
    let triggers = []
    let params = [database]
    if(triggerName){
      params.push(triggerName)
    }
    const results = await executeQuery(sql, params)
    triggers = await results.map(async result => {
      const name = result['TRIGGER_NAME']
      const content = result['CONTENT']
      const timing = result['ACTION_TIMING']
      const manipulation = result['EVENT_MANIPULATION']
      const table = result['EVENT_OBJECT_TABLE']
      const orientation = result['ACTION_ORIENTATION']
      let trigger = `CREATE TRIGGER ${name}
  ${timing} ${manipulation}
  ON ${table}
  FOR EACH ${orientation}
${content}`

      return { name: name, content: trigger }
    })

    return triggers
  } catch (err) {
    console.log(err)
    return []
  }
}

const getViews = async (viewName) => {
  let whereViewName = ''
  if(viewName){
    whereViewName = 'AND TABLE_NAME = ?'
  }
  const sql = `SELECT TABLE_SCHEMA, TABLE_NAME, VIEW_DEFINITION
    FROM INFORMATION_SCHEMA.VIEWS
    WHERE TABLE_SCHEMA = ?
    ${whereViewName}`

  try {
    let views = []
    let params = [database]
    if(whereViewName){
      params.push(viewName)
    }
    const results = await executeQuery(sql, params)
    views = await results.map(async result => {
      const name = result['TABLE_NAME']
      let content = result['VIEW_DEFINITION']
      content = content.split('select').join(`select\n `)
      content = content.split('SELECT').join(`SELECT\n `)
      content = content.split('from').join(`\nfrom`)
      content = content.split('FROM').join(`\nFROM`)
      content = content.split(',').join(`,\n  `)

      const options = ['right', 'left', 'inner', '', 'left outer', 'right outer']
      options.map((option) => {
        const token = `${option} join`.trim()
        const upperToken = token.toUpperCase()
        content = content.split(token).join(`\n${token}`)
        content = content.split(upperToken).join(`\n${upperToken}`)
      })
      let procedure = `CREATE OR REPLACE VIEW ${name} AS\n${content}`

      return { name: name, content: procedure }
    })

    return views
  } catch (error) {
    console.log(error)
    return []
  }
}

const getEvents = async (eventName) => {
  let whereEventName = ''
  if(eventName){
    whereEventName = 'AND EVENT_NAME = ?'
  }
  const sql = `SELECT * FROM INFORMATION_SCHEMA.EVENTS WHERE EVENT_SCHEMA = ? ${whereEventName}`

  try {
    let events = []
    let params = [database]
    if(eventName) {
      params.push([database])
    }
    const results = await executeQuery(sql, params)
    events = await results.map(async result => {
      const name = result['EVENT_NAME']
      const content = result['EVENT_DEFINITION']
      const executeAt = result['EXECUTE_AT']
      const interval = result['INTERVAL_VALUE']
      const intervalField = result['INTERVAL_FIELD']
      const comment = result['EVENT_COMMENT']
      const status = result['STATUS']
      const onCompletion = result['ON_COMPLETION']
      let schedule = `EVERY ${interval} ${intervalField}`
      if (executeAt !== null) {
        schedule = `AT ${executeAt}`
      }
      let event = `CREATE EVENT ${name}
  ON SCHEDULE ${schedule}
  ON COMPLETION ${onCompletion}
  COMMENT '${comment}'
  DO
${content};
ALTER EVENT ${name}
  ${status}`
      return { name: name, content: event }
    })

    return events
  } catch (error) {
    console.log(error)
    return []
  }
}

module.exports = {
  executeQuery,
  getProcedures,
  getFunctions,
  getTriggers,
  getViews,
  getTables,
  getEvents,
  end
}
