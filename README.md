# Planz 

Opentable reservation scraper

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
