const opentable = require('../')

let client = new opentable.Client()

client.search({
  latitude: '33.611746',
  longitude: '-117.7487',
  term: 'italian',
})
  .then(data => console.log(JSON.stringify(data, null, "  ")))
  .catch(err => console.log(err))

