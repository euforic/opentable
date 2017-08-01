
const exec = require('child_process').exec
const path = require("path")
const os = require("os")

const BIN_PATH = path.join(__dirname, `/bin/opentable-${os.platform()}`)

const DEFAULT_PROPS = {
  covers: 2, // number of people attending
  Limit: 200, // number of results to show
  enableSimpleCuisines: true,
  dateTime: new Date(), // date/time of reservation
  pageType: 0,
  latitude: '', // longitude of location to search
  longitude: '', // latitude of location to search
  term: '', // search term used to narrow search
  sort: 'Rating', // sort by [Distance, Popularity, Name, Rating]
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
  constructor(props) {
    this.baseURL = 'https://www.opentable.com/'
  }

  search(props) {
    let { agent } = props
    if (props.agent === '') {
      delete props.agent
    }

    let req = {
      url: makeSearchURL(Object.assign({}, DEFAULT_PROPS, props), this.baseURL),
      agent,
    }

    return new Promise((resolve, reject) => {
      exec(`${BIN_PATH} '${JSON.stringify(req)}'`, (error, stdout, stderr) => {
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

function makeSearchURL(props, baseURL) {
  const keys = Object.keys(props)
  const lastIndex = keys.length - 1

  let url = keys.reduce((sum, key, i) => {
    const val = (key === 'dateTime') ? formatDate(props[key]) : props[key]
    return sum += `${key}=${val}${i !== lastIndex ? '&' : ''}`
  }, `${baseURL}s/?`)

  return url
}

function pad(val) { 
  return (`0${val}`).slice(-2)
}

function formatDate(date) {
  return `${date.getFullYear()}-${pad(date.getMonth()+1)}-${pad(date.getDate())}+${pad(date.getHours())}:${pad(date.getMinutes())}`
}

