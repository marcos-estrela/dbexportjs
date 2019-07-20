const fs = require('fs')

const saveDbObjects = async (objectType, objectsList) => {
  const objects = await objectsList;
  try{
    objects.forEach(async objPromisse => {
      objPromisse.then(async obj => {
        await write(`./export/${objectType}`,`${obj.name}.sql`, obj.content)
      })
    })
  } catch(err) {
    console.log(err)
  }
}

const makeDir = (dirName) => {
  try{
    if(!fs.existsSync(dirName)){
      fs.mkdirSync(dirName, {recursive: true});
    }
    return true;
  }catch(err) {
    return false;
  }
}

const write = (dirName, fileName, fileContent) => {
  const path = `${dirName}/${fileName}`

  if (!fs.existsSync(dirName)){
    const dirParts = dirName.replace('./', '').split('/')
    let fullPath = '.'
    dirParts.forEach(part => {
      fullPath = `${fullPath}/${part}`
      makeDir(fullPath)
    })
  }

  fs.writeFile(path, fileContent, function (err) {
    if (err) {
      return console.log(err);
    }

    console.log(`The file ${fileName} was saved!`);
  });
}

module.exports = {
  saveDbObjects,
  write
}
