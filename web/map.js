var map = L.map('map');
map.on('click', onMapClick);
map.on('move', onMapMove);
map.on('zoom', onMapZoom);
map.on('load', onMapLoad);
map.on('error', function(e) {
  console.log('Error loading file: ' + e.err);
});

// Set map at location and zoom value
// Commented as the is update on gpx loaded
// map.setView([45.742, 6.523], 12);

L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map); 



const options = {
    async: true,
    // polyline_options: { color: 'red' },
    markers: {

        startIcon: "imgs/placeholder.png", // impossble d'utiliser L.icon
        // endIcon: "",

        wptTypeIcons: {
            "refuge": "imgs/person-shelter.svg",
            "camping": "imgs/tent.svg",
            "ville": "imgs/building.svg",
            "river": "imgs/water.svg",
            "col": "imgs/mountain-sun.svg",
            "sommet": "imgs/mountain.svg",
            "vue": "imgs/eye.svg",
        },

        // Marche pas
        // pointMatchers: [
        //     {
        //         regex: /refuge/,
        //         icon: "imgs/location-dot.svg",
        //     },
        // ],

    },
    marker_options: {
        iconSize: [15, 15],
        iconAnchor: [15, 15],
    },
    gpx_options: {
        parseElements: ["track", "waypoint"],
        // joinTrackSegments: false
    },
};

// =====
// afficher les noms de trk
// =====

// function loadGPX(gpx) {
//   new L.GPX(gpx,{
//       async: true,
//       marker_options: {
//         startIcongpxpath:  'pin-icon-start.png',
//         endIcongpxpath:  false,
//         shadowgpxpath:   false,
//       },
//       gpx_options: {
//           joinTrackSegments: false
//       },
//   }).on('loaded', function(e) {
//     var track = e.target;
//     // control.addOverlay(track,track.get_name()) // adds a tickbox in the control panel of the map
//     console.log(track.get_name())              // to see in the browser console
//     console.log('distance = '+track.get_distance())
//     console.log('max ele =' + track.get_elevation_max())
//   }).on('addpoint', (e) => {
//       e.point.bindPopup(e.target.get_name());
//   }).addTo(map)
// }

// var req = new window.XMLHttpRequest();
// req.open('GET', gpxpath, true);
// try {
//   req.overrideMimeType('text/xml');
// } catch(e) {}
// req.onreadystatechange = function() {
//   if (req.readyState != 4) return;
//   if(req.status == 200) {
//     var xml = req.responseXML;
//     var tracks = xml.getElementsByTagName('trk');
//     var serializer = new XMLSerializer();
//     for (var i = 0; i < tracks.length; i++) {
//       loadGPX(serializer.serializeToString(tracks[i]));
//     }
//   }
// };
// req.send(null);

// =====

// =====
// ... pas les titre
// =====

const gpx_selection = document.getElementById("gpx_selection")
var gpxpath = gpx_selection.selectedOptions[0].value;
console.log(gpx_selection);
var distance, ascent, descent
var gpx;

function updateGpxInfo(distance, ascent, descent){
    document.getElementById("gpx-distance").innerHTML = Math.floor(distance / 1000)
    document.getElementById("gpx-ascent").innerHTML = Math.floor(ascent)
    document.getElementById("gpx-descent").innerHTML = Math.floor(descent)
}

function plotElevation(d){
    var trace1 = {
        x: d.map(e => e[0]),
        y: d.map(e => e[1]),
        type: 'scatter'
    };

    var data = [trace1];
    Plotly.newPlot('plot', data);
}


function updateGpx(){
    gpxpath = gpx_selection.selectedOptions[0].value;

    if (gpx !== undefined) {
        gpx.remove();
    }
    gpx = new L.GPX(gpxpath, options);
    gpx.on('loaded', (e) => {
        map.fitBounds(e.target.getBounds());

        updateGpxInfo(
            e.target.get_distance(),
            e.target.get_elevation_gain(),
            e.target.get_elevation_loss()
        )

        plotElevation(e.target.get_elevation_data())
    });

    // gpx.on('addpoint', function(e) {
    //     if (e.point_type === 'start' || e.point_type === 'end') {
    //         var text = e.target._info.name;
    //         text += '<br/>L&auml;nge: ' + e.target._info.length.toFixed(2) + ' m';
    //         text += '<br/>Anstieg: ' + e.target._info.elevation.gain.toFixed(2) + ' m';
    //         if(e.target._info.duration.start !== null)
    //               text += '<br/>Datum: ' + e.target._info.duration.start.toLocaleDateString();
    //         e.point.bindPopup(text);

    //         // this.removeLayer(e.point);
    //     }
    // });

    gpx.addTo(map);
    console.log(gpx);
};

updateGpx();
gpx_selection.addEventListener('change', updateGpx);

// ===


function getMapSWNEBounds(){
    bounds = map.getBounds()
    s = "".concat([ 
        bounds.getSouth(), 
        bounds.getWest(), 
        bounds.getNorth(),
        bounds.getEast(),
    ])
    return s
}

function onMapLoad(e) {
    document.getElementById("map-bounds").innerHTML = getMapSWNEBounds()
    document.getElementById("map-center").innerHTML = map.getCenter()
    document.getElementById("map-zoom").innerHTML = map.getZoom()
}

function onMapClick(e) {
    // console.log("You clicked the map at " + e.latlng);
    document.getElementById("click-position").innerHTML = e.latlng;
}

function onMapMove(e) {
    document.getElementById("map-bounds").innerHTML = getMapSWNEBounds()
    document.getElementById("map-center").innerHTML = map.getCenter()
}

function onMapZoom(e) {
    document.getElementById("map-zoom").innerHTML = map.getZoom()
}


