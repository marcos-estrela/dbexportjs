const pwd = process.cwd()
const ENV_FILE = process.env.NODE_ENV ? `${pwd}/.env.${process.env.NODE_ENV}` : `${pwd}/.env.local`

require('dotenv').config({ path: ENV_FILE })
