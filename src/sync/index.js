require('../config')

const adapters = require('../adapters')
const fs = require('fs')

class DbSync {
  constructor(adapter) {
    this.adapter = adapters[adapter]

    if(!this.adapter){
      throw new Error(`Adapter ${adapter} not found!`);
    }
  }

  async commitChanges(fileName) {
    try{
      const content = await this.getContentFromFile(fileName)
      if(content !== ''){
        content = this.addDropQueryIfNotExists(content);
        this.commit(content)
      }
    }catch(err){
      throw new Error(err)
    }
  }

  async getContentFromFile(fileName) {
    return new Promise(function(resolve, reject) {
        fs.readFile(fileName, "utf8", (err, content) => {
          if(err) reject(err)
          resolve(content)
        })
    })
  }

  async addDropQueryIfNotExists(content) {
    let lines = content.split('\n')
    let header = lines[0]
    let coluns = header.split('(')[0].split(' ')
    let objParts = {
        actionPos: 0,
        typePos: 1,
        namePos: 2
    }
    let newHeader = ''

    if(coluns[objParts.actionPos].toLowerCase() !== 'drop') {
        let objType = coluns[objParts.typePos]
        let objName = coluns[objParts.namePos]
        if(!header.includes('view')){
            newHeader = `DROP ${objType} IF EXISTS ${objName};\n\n`
        }
    }

    content = `${newHeader}${content};\n`
    return content
  }

  async commit(query) {
    try{
      return this.adapter.executeQuery(query)
    }catch(err){
      console.log(err)
      return null
    }
  }
}

module.exports = DbSync
