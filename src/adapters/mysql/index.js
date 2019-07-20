const fs = require('fs')
const connection = require('./connection')

const database = process.env.DATABASE

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
  const sql = `SELECT SPECIFIC_NAME, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER
    FROM INFORMATION_SCHEMA.PARAMETERS
    WHERE
    SPECIFIC_SCHEMA = ?
    AND SPECIFIC_NAME = ?
    AND ROUTINE_TYPE = ?
    ORDER BY SPECIFIC_NAME, ORDINAL_POSITION ASC`
  try {
    const results = await executeQuery(sql, [database, objectName, objectType])
    return results
  } catch (err) {
    console.log(err)
    return []
  }
}

const getTables = async () => {
  const sql = `SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = ? AND TABLE_SCHEMA = ?`

  try {
    let tables = []
    const results = await executeQuery(sql, ['BASE TABLE', database])
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

const getFunctions = async () => {
  const sql = `SELECT ROUTINE_NAME, ROUTINE_TYPE, ROUTINE_DEFINITION, ROUTINE_COMMENT, DTD_IDENTIFIER
    FROM INFORMATION_SCHEMA.ROUTINES
    WHERE
    ROUTINE_SCHEMA = ?
    AND ROUTINE_TYPE = 'FUNCTION'`
  try {
    let functions = []
    const results = await executeQuery(sql, [database])
    functions = await results.map(async result => {
      const name = result['ROUTINE_NAME']
      const charSetName = result['CHARACTER_SET_NAME']
      const charSet = charSetName ? ` CHARSET ${charSetName}` : ''
      let parametersList = []
      let parameters = await getParameters(name, 'FUNCTION')
      parametersList = await parameters.map(param => {
        return `${param['PARAMETER_NAME']} ${param['DTD_IDENTIFIER']}`
      })
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

const getProcedures = async () => {
  const sql = `SELECT ROUTINE_NAME, ROUTINE_TYPE, ROUTINE_DEFINITION, ROUTINE_COMMENT
    FROM INFORMATION_SCHEMA.ROUTINES
    WHERE
    ROUTINE_SCHEMA = ?
    AND ROUTINE_TYPE = 'PROCEDURE'`
  try {
    let procedures = []
    const results = await executeQuery(sql, [database])
    procedures = await results.map(async result => {
      const name = result['ROUTINE_NAME']
      const comment = result['ROUTINE_COMMENT']
      let parametersList = []
      let parameters = await getParameters(name, 'PROCEDURE')
      parametersList = await parameters.map(param => {
        return `${param['PARAMETER_MODE']} ${param['PARAMETER_NAME']} ${param['DTD_IDENTIFIER']}`
      })
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

const getTriggers = async () => {
  const sql = `SELECT TRIGGER_NAME, ACTION_STATEMENT AS CONTENT, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORIENTATION
    FROM INFORMATION_SCHEMA.TRIGGERS
    WHERE
    TRIGGER_SCHEMA = ?`

  try {
    let triggers = []
    const results = await executeQuery(sql, [database])
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

const getViews = async () => {
  const sql = `SELECT TABLE_SCHEMA, TABLE_NAME, VIEW_DEFINITION
    FROM INFORMATION_SCHEMA.VIEWS
    WHERE TABLE_SCHEMA = ?`

  try {
    let views = []
    const results = await executeQuery(sql, [database])
    views = await results.map(async result => {
      const name = result['TABLE_NAME']
      const content = result['VIEW_DEFINITION']
      let procedure = `CREATE OR REPLACE VIEW ${name} AS ${content}`

      return { name: name, content: procedure }
    })

    return views
  } catch (error) {
    console.log(error)
    return []
  }
}

const getEvents = async () => {
  const sql = `SELECT * FROM INFORMATION_SCHEMA.EVENTS WHERE EVENT_SCHEMA = ?`

  try {
    let events = []
    const results = await executeQuery(sql, [database])
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
  getProcedures,
  getFunctions,
  getTriggers,
  getViews,
  getTables,
  getEvents,
  end
}
