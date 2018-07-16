const path = require('path')
const express = require('express')
const exphbs = require('express-handlebars')
var Client = require('node-rest-client').Client;

const port=3000

var client = new Client();
const app = express()

app.engine('.hbs', exphbs({
  defaultLayout: 'main',
  extname: '.hbs',
  layoutsDir: path.join(__dirname, 'views/layouts')
}))
app.set('view engine', '.hbs')
app.set('views', path.join(__dirname, 'views'))

app.use(express.static(path.join(__dirname, "static")));

app.get('/devicedata/alerts', (req, resp) => {
  
  client.get("http://localhost:8888/devicedata?device=" + req.query.id, function (data , response) {
    
    var result = {};
    for (var i=0;i < data.length;i++) {
      var entry = data[i];
      var res_entry = {};
      res_entry.time = new Date(entry.blocktime*1000);
      if (entry.data.json.rule != undefined) {
        // this is an alert
        res_entry.type = "alert";
        res_entry.rule = entry.data.json.rule;
        res_entry.data = entry.data.json.data;
      } else {
        // This is a clear alert
        res_entry.type = "clear";
        res_entry.resolution = entry.data.json.Resolution;
        res_entry.alertid = entry.data.json.AlertTxid;
      }
      result[entry.txid] = res_entry;
    }
    resp.json(result)
 });   
})

app.get('/devicedata/monitor', (req, resp) => {

  client.get("http://localhost:8888/devicedata?device=" + req.query.id + "&type=monitor", function (data , response) {
    
    var result = {};
    for (var i=0;i < data.length;i++) {
      var entry = data[i];
      var res_entry = {};
      res_entry.time = new Date(entry.blocktime*1000);
      res_entry.temperature = entry.data.json.temperature;
      res_entry.humidity = entry.data.json.humidity;
      result[entry.txid] = res_entry;
    }
    resp.json(result)
  });   

})

app.get('/device', (request, response) => {
  response.render('home', {
    devid: request.query.id
  })
})

app.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err)
  }

  console.log(`server is listening on ${port}`)
})

