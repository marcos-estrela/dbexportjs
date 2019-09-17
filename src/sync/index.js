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
