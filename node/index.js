
const exec = require('child_process').exec
const path = require("path")
const os = require("os")

const BIN_PATH = path.join(__dirname, `/bin/opentable-${os.platform()}`)

const DEFAULT_PROPS = {
  people: 2, // number of people attending
  limit: 200, // number of results to show
  dateTime: new Date(), // date/time of reservation
  latitude: '', // longitude of location to search
  longitude: '', // latitude of location to search
  term: '', // search term used to narrow search
  sort: 'Rating', // sort by [Distance, Popularity, Name, Rating]
  opts: {},
}

const flags = {
  people : '--people',
  //  limit: '--limit',
  dateTime: '--date_time',
  longitude: '--longitude',
  latitude: '--latitude',
  term: '--term',
  sort: '--sort',
  agent: '--agent',
}


// Example usage
//
// let client = new Client()
// client.search({
//   latitude: '33.611746',
//   longitude: '-117.7487',
//   term: 'italian',
// }).then(data => console.log(data)).catch(err => console.log(err))
// 

class Client {
  _formatArgs(props) {
    let opts = Object.assign({}, DEFAULT_PROPS, props)
    opts.dateTime = formatDate(opts.dateTime)

    let args = Object.keys(opts).reduce((str, itm) => {
      if (!flags[itm]) { return str }
      return `${str} ${flags[itm]} ${opts[itm]} `
    }, '')

    if (!opts.opts || Object.keys(opts.opts).length === 0){
      return args
    }

    args += Object.keys(opts.opts).reduce((str, itm) => {
      return `${str}${itm}=${opts.opts[itm]},`
    }, '--opts ')
    return args
  }

  search(props) {
    return new Promise((resolve, reject) => {
      console.log(`${BIN_PATH} client search --direct ${this._formatArgs(props)}`)
      exec(`${BIN_PATH} client search --direct ${this._formatArgs(props)}`, (error, stdout, stderr) => {
        if (error) {
          reject(error)
          return
        }

        if (stderr !== '') {
          reject(new Error(stderr))
          return
        }

        try {
          let data = JSON.parse(stdout)
          resolve(data)
        } catch(error) {
          reject(new Error(error))
        }
      })
    })
  }
}

exports.Client = Client

function pad(val) { 
  return (`0${val}`).slice(-2)
}

function formatDate(date) {
  return `${date.getFullYear()}-${pad(date.getMonth()+1)}-${pad(date.getDate())}+${pad(date.getHours())}:${pad(date.getMinutes())}`
}

