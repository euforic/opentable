const opentable = require('../node')

let x = new opentable.Client()
x.search({
  latitude: '33.611746',
  longitude: '-117.7487',
  term: 'italian',
}).then(data => console.log(data)).catch(err => console.log(err))

