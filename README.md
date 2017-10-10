# Opentable 

Opentable reservation scraper in golang with nodejs client and gRPC server / client.
This code base is meant to serve as an example for creating a simple scraper, gRPC service,
and Golang <=> Node.js communication. Odds are this will not be maintained. :-(

**Personal use only! May violate Opentable's [T&C #13](https://www.opentable.com/legal/terms-and-conditions).**

[![GoDoc](https://godoc.org/github.com/euforic/opentable?status.svg)](https://godoc.org/github.com/euforic/opentable)
![](https://img.shields.io/badge/license-MIT-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/euforic/opentable)](https://goreportcard.com/report/github.com/euforic/opentable)

## CLI Usage

```bash
$ opentable client search --direct  --people 2  --date_time 2017-08-11+15:30  --latitude 33.611746  --longitude -117.7487  --term italian  --sort Rating
```

## Nodejs Usage

### Request

```js
let client = new opentable.Client()

client.search({
  covers: 2, // number of people attending
  Limit: 200, // number of results to show
  dateTime: new Date(), // date/time of reservation
  term: 'italian', // search term used to narrow search
  sort: 'Rating', // sort by [Distance, Popularity, Name, Rating]
  latitude: '33.611746', // latitude of location to search
  longitude: '-117.7487', // longitude of location to search
}).then(data => {
  console.log(JSON.stringify(data, null, "  "))
}).catch(err => {
  console.log(err)
})
```

### Response

```js
{
  "350182": {
    "ID": "350182",
    "Name": "Alessa by Chef Pirozzi",
    "URL": "https://opentable.com/restaurant/profile/350182?p=2&sd=2017-08-02%2014%3A30",
    "Recommended": "97%",
    "Reservations": [
      {
        "Time": "2017-08-02T14:30:00Z",
        "URL": "https://opentable.com/book/validate?rid=350182&d=2017-08-02 14:30&p=2&pt=100&ss=0&sd=2017-08-02 14:30&pofids=&hash=490905459"
      },
      //...
    ]
  },
  //...
}
