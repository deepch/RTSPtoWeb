let suuid = $('#suuid').val();

let config = {
  iceServers: [{
    urls: ["stun:stun.l.google.com:19302"]
  }]
};

const pc = new RTCPeerConnection(config);
pc.onnegotiationneeded = handleNegotiationNeededEvent;

let log = msg => {
  document.getElementById('div').innerHTML += msg + '<br>'
}

pc.ontrack = function(event) {
  log(event.streams.length + ' track is delivered')
  var el = document.createElement(event.track.kind)
  el.srcObject = event.streams[0]
  el.muted = true
  el.autoplay = true
  el.controls = true
  el.width = 600
  document.getElementById('remoteVideos').appendChild(el)
}

pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)



async function handleNegotiationNeededEvent() {
  let offer = await pc.createOffer();
  await pc.setLocalDescription(offer);
  getRemoteSdp();
}

$(document).ready(function() {
 // $('#' + suuid).addClass('active');
  getCodecInfo();
});


function getCodecInfo() {
  $.get("/codec/demo2" + suuid, function(data) {
    try {
      data = JSON.parse(data);
      if (data.length > 1) {
        log('add audio Transceiver')
        pc.addTransceiver('audio', {
          'direction': 'sendrecv'
        })
      }
    } catch (e) {
      console.log(e);
    } finally {

      log('add video Transceiver')
      pc.addTransceiver('video', {
        'direction': 'sendrecv'
      });
      //send ping becouse PION not handle RTCSessionDescription.close()
      sendChannel = pc.createDataChannel('foo');
      sendChannel.onclose = () => console.log('sendChannel has closed');
      sendChannel.onopen = () => {
        console.log('sendChannel has opened');
        sendChannel.send('ping');
        setInterval(() => {
          sendChannel.send('ping');
        }, 1000)
      }
      sendChannel.onmessage = e => log(`Message from DataChannel '${sendChannel.label}' payload '${e.data}'`);
    }
  });
}

let sendChannel = null;

function getRemoteSdp() {
  $.post("/stream/demo2/webrtc", {
    suuid: suuid,
    data: btoa(pc.localDescription.sdp)
  }, function(data) {
    try {

      pc.setRemoteDescription(new RTCSessionDescription({
        type: 'answer',
        sdp: atob(data)
      }))



    } catch (e) {
      console.warn(e);
    }

  });
}


if (!Uint8Array.prototype.slice) {
  Object.defineProperty(Uint8Array.prototype, 'slice', {
    value: function(begin, end) {
      return new Uint8Array(Array.prototype.slice.call(this, begin, end));
    }
  });
}

var verbose = true;
var streamingStarted = false;
var ms = new MediaSource();
var queue = [];
var ws;

function pushPacket(arr) {
  var view = new Uint8Array(arr);
  if (verbose) {
    console.log("got", arr.byteLength, "bytes.  Values=", view[0], view[1], view[2], view[3], view[4]);
  }
  data = arr;
  if (!streamingStarted) {
    sourceBuffer.appendBuffer(data);
    streamingStarted = true;
    return;
  }
  queue.push(data);
  if (verbose) {
    console.log("queue push:", queue.length);
  }
  if (!sourceBuffer.updating) {
    loadPacket();
  }
}

function loadPacket() {
  if (!sourceBuffer.updating) {
    if (queue.length > 0) {
      inp = queue.shift();
      if (verbose) {
        console.log("queue PULL:", queue.length);
      }
      var view = new Uint8Array(inp);
      if (verbose) {
        console.log("writing buffer with", view[0], view[1], view[2], view[3], view[4]);
      }
      sourceBuffer.appendBuffer(inp);
    } else {
      streamingStarted = false;
    }
  }
}

var potocol = 'ws';
if (location.protocol.indexOf('s') >= 0) {
  potocol = 'wss';
}

function opened() {
  var inputVal = $('#suuid').val();
  var port = $('#port').val();
  ws = new WebSocket(potocol + "://127.0.0.1"+port+"/stream/demo2/mse?uuid=demo2");
  ws.binaryType = "arraybuffer";
  ws.onopen = function(event) {
    console.log('Connect');
  }
  ws.onmessage = function(event) {
    var data = new Uint8Array(event.data);
    if (data[0] == 9) {
      decoded_arr=data.slice(1);
      if (window.TextDecoder) {
        mimeCodec = new TextDecoder("utf-8").decode(decoded_arr);
      } else {
        mimeCodec = Utf8ArrayToStr(decoded_arr);
      }
      if (verbose) {
        console.log('first packet with codec data: ' + mimeCodec);
      }
      sourceBuffer = ms.addSourceBuffer('video/mp4; codecs="' + mimeCodec + '"');
      sourceBuffer.mode = "segments"
      sourceBuffer.addEventListener("updateend", loadPacket);
    } else {
      pushPacket(event.data);
    }
  };
}
var livestream = document.getElementById('livestream');

function Utf8ArrayToStr(array) {
  var out, i, len, c;
  var char2, char3;
  out = "";
  len = array.length;
  i = 0;
  while (i < len) {
    c = array[i++];
    switch (c >> 4) {
      case 7:
        out += String.fromCharCode(c);
        break;
      case 13:
        char2 = array[i++];
        out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));
        break;
      case 14:
        char2 = array[i++];
        char3 = array[i++];
        out += String.fromCharCode(((c & 0x0F) << 12) |
          ((char2 & 0x3F) << 6) |
          ((char3 & 0x3F) << 0));
        break;
    }
  }
  return out;
}

function startup() {
  ms.addEventListener('sourceopen', opened, false);
  livestream.src = window.URL.createObjectURL(ms);
}

$(document).ready(function() {
startup();
  //var suuid = $('#suuid').val();
 // $('#'+suuid).addClass('active');
});
