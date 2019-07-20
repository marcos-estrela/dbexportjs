const fs = require('fs')

const loadJsonConfig = (jsonFile) => {
  const rawdata = fs.readFileSync(jsonFile)
  const config = JSON.parse(rawdata)
  process.env = {...process.env, ...config}
}

const loadEnvConfig = (envFile) => {
  require('dotenv').config({ path: envFile })
}

const getEnvFile = () => {
  const pwd = process.cwd()
  const envFile = process.env.NODE_ENV ? `${pwd}/.env.${process.env.NODE_ENV}` : `${pwd}/.env.local`
  return envFile
}

const loadConfig = () => {
  const jsonFile = './config.json'
  const envFile = getEnvFile()

  if(fs.existsSync(jsonFile)){
    loadJsonConfig(jsonFile)
  }else if(fs.existsSync(envFile)){
    loadEnvConfig(envFile)
  }else{
    throw new Error(`One of this files should be exists: ${file} | ${ENV_FILE}`)
  }
}

try{
  loadConfig()
}catch(err) {
  console.error(err)
}
