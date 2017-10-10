const opentable = require('../')

let client = new opentable.Client()

client.search({
  latitude: '33.611746',
  longitude: '-117.7487',
  term: 'italian',
  opts: {
    Limit: 3,
  },
})
  .then(data => console.log(JSON.stringify(data, null, "  ")))
  .catch(err => console.log(err))

