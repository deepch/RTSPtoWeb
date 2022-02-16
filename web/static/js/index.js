$(document).ready(() => {
  localImages();
  if (localStorage.getItem('defaultPlayer') != null) {
    $('input[name=defaultPlayer]').val([localStorage.getItem('defaultPlayer')]);
  }
})
$('input[name=defaultPlayer]').on('change', function() {
  localStorage.setItem('defaultPlayer', $(this).val());
})

var activeStream = null;

function showAddStream(streamName, streamUrl) {
  streamName = streamName || '';
  streamUrl = streamUrl || '';
  Swal.fire({
    title: 'Add stream',
    html: '<form class="text-left"> ' +
      '<div class="form-group">' +
      '<label>Name</label>' +
      '<input type="text" class="form-control" id="stream-name">' +
      '<small class="form-text text-muted"></small>' +
      '</div>' +
      '<div class="form-group">' +
      '  <label>URL</label>' +
      '  <input type="text" class="form-control" id="stream-url">' +
      '  </div>' +
      '<div class="form-group form-check">' +
      '<input type="checkbox" class="form-check-input" id="stream-ondemand">' +
      '<label class="form-check-label">ondemand</label>' +
      '</div>' +
      '</form>',
    focusConfirm: true,
    showCancelButton: true,
    preConfirm: () => {
      var uuid = randomUuid(),
        name = $('#stream-name').val(),
        url = $('#stream-url').val(),
        ondemand = $('#stream-ondemand').val();
      if (!validURL(url)) {
        Swal.fire({
          icon: 'error',
          title: 'Oops...',
          text: 'wrong url',
          confirmButtonText: 'return back',
          preConfirm: () => {
            showAddStream(name, url)
          }
        })
      } else {
        goRequest('add', uuid, {
          name: name,
          url: url,
          ondemand: ondemand
        });


      }
    }
  })

}

function showEditStream(uuid) {
  console.log(streams[uuid]);
}

function deleteStream(uuid) {
  activeStream = uuid;
  Swal.fire({
    title: 'Are you sure?',
    text: "Do you want delete stream " + streams[uuid].name + " ?",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#3085d6',
    cancelButtonColor: '#d33',
    confirmButtonText: 'Yes, delete it!'
  }).then((result) => {
    if (result.value) {
      goRequest('delete', uuid)

    }
  })
}

function renewStreamlist() {
  goRequest('streams');
}

function goRequest(method, uuid, data) {
  data = data || null;
  uuid = uuid || null;
  var path = '';
  var type = 'GET';
  switch (method) {
    case 'add':
      path = '/stream/' + uuid + '/add';
      type = 'POST';
      break;
    case 'edit':
      path = '/stream/' + uuid + '/edit';
      type = 'POST';
      break;
    case 'delete':
      path = '/stream/' + uuid + '/delete';
      break;
    case 'reload':
      path = '/stream/' + uuid + '/reload';
      break;
    case 'info':
      path = '/stream/' + uuid + '/info';
      break;
    case 'streams':
      path = '/streams';
      break;
    default:
      path = '';
      type = 'GET';
  }
  if (path == '') {
    Swal.fire({
      icon: 'error',
      title: 'Oops...',
      text: 'It`s goRequest function mistake',
      confirmButtonText: 'Close',

    })
    return;
  }
  var ajaxParam = {
    url: path,
    type: type,
    dataType: 'json',
    beforeSend: function(xhr) {
      xhr.setRequestHeader("Authorization", "Basic " + btoa("demo:demo"));
    },
    success: function(response) {
      goRequestHandle(method, response, uuid);
    },
    error: function(e) {
      console.log(e);
    }
  };
  if (data != null) {
    ajaxParam.data = JSON.stringify(data);
  }
  $.ajax(ajaxParam);
}

function goRequestHandle(method, response, uuid) {
  switch (method) {
    case 'add':

      if (response.status == 1) {
        renewStreamlist();
        Swal.fire(
          'Added!',
          'Your stream has been Added.',
          'success'
        );

      } else {
        Swal.fire({
          icon: 'error',
          title: 'Oops...',
          text: 'Same mistake issset',
        })
      }

      break;
    case 'edit':
      if (response.status == 1) {
        renewStreamlist();
        Swal.fire(
          'Success!',
          'Your stream has been modified.',
          'success'
        );
      } else {
        Swal.fire({
          icon: 'error',
          title: 'Oops...',
          text: 'Same mistake issset',
        })
      }
      break;
    case 'delete':

      if (response.status == 1) {
        $('#' + uuid).remove();
        delete(streams[uuid]);
        $('#stream-counter').html(Object.keys(streams).length);
        Swal.fire(
          'Deleted!',
          'Your stream has been deleted.',
          'success'
        )
      }

      break;
    case 'reload':

      break;
    case 'info':

      break;
    case 'streams':
      if (response.status == 1) {
        streams = response.payload;
        $('#stream-counter').html(Object.keys(streams).length);
        if (Object.keys(streams).length > 0) {

          $.each(streams, function(uuid, param) {
            if ($('#' + uuid).length == 0) {
              $('.streams').append(streamHtmlTemplate(uuid, param.name));
            }
          })
        }
      }

      break;
    default:

  }

}

function getImageBase64(videoEl){
    const canvas = document.createElement("canvas");
    canvas.width = videoEl.videoWidth;
    canvas.height = videoEl.videoHeight;
    canvas.getContext('2d').drawImage(videoEl, 0, 0, canvas.width, canvas.height);
    const dataURL = canvas.toDataURL();
    canvas.remove();
    return dataURL;
}

function downloadBase64Image(base64Data){
    const a = document.createElement("a");
    a.href = base64Data;
    a.download = "screenshot.png";
    a.click();
    a.remove();
}


function makePic(video_element, uuid, chan) {
  if (typeof(video_element) === "undefined") {
    video_element = $("#videoPlayer")[0];
  }
  ratio = video_element.videoWidth / video_element.videoHeight;
  w = 400;
  h = parseInt(w / ratio, 10);
  $('#canvas')[0].width = w;
  $('#canvas')[0].height = h;
  $('#canvas')[0].getContext('2d').fillRect(0, 0, w, h);
  $('#canvas')[0].getContext('2d').drawImage(video_element, 0, 0, w, h);
  var imageData = $('#canvas')[0].toDataURL();
  var images = localStorage.getItem('imagesNew');
  if (images != null) {
    images = JSON.parse(images);
  } else {
    images = {};
  }
  var uid = $('#uuid').val();
  if (!!uuid) {
    uid = uuid;
  }

  var channel = $('#channel').val() || chan || 0;
  if (typeof(images[uid]) === "undefined") {
    images[uid] = {};
  }
  images[uid][channel] = imageData;
  localStorage.setItem('imagesNew', JSON.stringify(images));
  $('#' + uid).find('.stream-img[channel="' + channel + '"]').attr('src', imageData);
}

function localImages() {
  var images = localStorage.getItem('imagesNew');
  if (images != null) {
    images = JSON.parse(images);
    $.each(images, function(k, v) {
      $.each(v, function(channel, img) {
        $('#' + k).find('.stream-img[channel="' + channel + '"]').attr('src', img);
      })

    });
  }
}

function clearLocalImg() {
  localStorage.setItem('imagesNew', '{}');
}

function streamHtmlTemplate(uuid, name) {
  return '<div class="item" id="' + uuid + '">' +
    '<div class="stream">' +
    '<div class="thumbs" onclick="rtspPlayer.livePlayer(0, \'' + uuid + '\')">' +
    '<img src="../static/img/noimage.svg" alt="" class="stream-img">' +
    '</div>' +
    '<div class="text">' +
    '<h5>' + name + '</h5>' +
    '<p>property</p>' +
    '<div class="input-group-prepend dropleft text-muted">' +
    '<a class="btn" data-toggle="dropdown" >' +
    '<i class="fas fa-ellipsis-v"></i>' +
    '</a>' +
    '<div class="dropdown-menu">' +
    '<a class="dropdown-item" onclick="rtspPlayer.livePlayer(\'hls\', \'' + uuid + '\')" href="#">Play HLS</a>' +
    '<a class="dropdown-item" onclick="rtspPlayer.livePlayer(\'mse\', \'' + uuid + '\')" href="#">Play MSE</a>' +
    '<a class="dropdown-item" onclick="rtspPlayer.livePlayer(\'webrtc\', \'' + uuid + '\')" href="#">Play WebRTC</a>' +
    '<div class="dropdown-divider"></div>' +
    '<a class="dropdown-item" onclick="showEditStream(\'' + uuid + '\')" href="#">Edit</a>' +
    '<a class="dropdown-item" onclick="deleteStream(\'' + uuid + '\')" href="#">Delete</a>' +
    '</div>' +
    '</div>' +
    '</div>' +
    '</div>' +
    '</div>';
}

function randomUuid() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

function validURL(url) {
  //TODO: fix it
  return true;
}

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

function browserDetector() {
  var Browser;
  var ua = self.navigator.userAgent.toLowerCase();
  var match =
    /(edge)\/([\w.]+)/.exec(ua) ||
    /(opr)[\/]([\w.]+)/.exec(ua) ||
    /(chrome)[ \/]([\w.]+)/.exec(ua) ||
    /(iemobile)[\/]([\w.]+)/.exec(ua) ||
    /(version)(applewebkit)[ \/]([\w.]+).*(safari)[ \/]([\w.]+)/.exec(ua) ||
    /(webkit)[ \/]([\w.]+).*(version)[ \/]([\w.]+).*(safari)[ \/]([\w.]+)/.exec(
      ua
    ) ||
    /(webkit)[ \/]([\w.]+)/.exec(ua) ||
    /(opera)(?:.*version|)[ \/]([\w.]+)/.exec(ua) ||
    /(msie) ([\w.]+)/.exec(ua) ||
    (ua.indexOf("trident") >= 0 && /(rv)(?::| )([\w.]+)/.exec(ua)) ||
    (ua.indexOf("compatible") < 0 && /(firefox)[ \/]([\w.]+)/.exec(ua)) || [];
  var platform_match =
    /(ipad)/.exec(ua) ||
    /(ipod)/.exec(ua) ||
    /(windows phone)/.exec(ua) ||
    /(iphone)/.exec(ua) ||
    /(kindle)/.exec(ua) ||
    /(android)/.exec(ua) ||
    /(windows)/.exec(ua) ||
    /(mac)/.exec(ua) ||
    /(linux)/.exec(ua) ||
    /(cros)/.exec(ua) || [];
  var matched = {
    browser: match[5] || match[3] || match[1] || "",
    version: match[2] || match[4] || "0",
    majorVersion: match[4] || match[2] || "0",
    platform: platform_match[0] || ""
  };
  var browser = {};

  if (matched.browser) {
    browser[matched.browser] = true;
    var versionArray = matched.majorVersion.split(".");
    browser.version = {
      major: parseInt(matched.majorVersion, 10),
      string: matched.version
    };

    if (versionArray.length > 1) {
      browser.version.minor = parseInt(versionArray[1], 10);
    }

    if (versionArray.length > 2) {
      browser.version.build = parseInt(versionArray[2], 10);
    }
  }

  if (matched.platform) {
    browser[matched.platform] = true;
  }

  if (browser.chrome || browser.opr || browser.safari) {
    browser.webkit = true;
  } // MSIE. IE11 has 'rv' identifer

  if (browser.rv || browser.iemobile) {
    if (browser.rv) {
      delete browser.rv;
    }

    var msie = "msie";
    matched.browser = msie;
    browser[msie] = true;
  } // Microsoft Edge

  if (browser.edge) {
    delete browser.edge;
    var msedge = "msedge";
    matched.browser = msedge;
    browser[msedge] = true;
  } // Opera 15+

  if (browser.opr) {
    var opera = "opera";
    matched.browser = opera;
    browser[opera] = true;
  } // Stock android browsers are marked as Safari

  if (browser.safari && browser.android) {
    var android = "android";
    matched.browser = android;
    browser[android] = true;
  }

  browser.name = matched.browser;
  browser.platform = matched.platform;


  return browser;
}

function addChannel() {
  $('#streams-form-wrapper').append(chanellTemplate());
}

function chanellTemplate() {
  let random = Math.ceil(Math.random() * 1000);
  let html = `
    <div class="col-12">
      <div class="card card-secondary">
        <div class="card-header">
          <h3 class="card-title">Sub channel<small> parameters</small></h3>
          <div class="card-tools">
          <button type="button" class="btn btn-tool" onclick="removeChannelDiv(this)"><i class="fas fa-times"></i></button>
          </div>
        </div>
          <div class="card-body">
          <form class="stream-form">
            <div class="form-group">
              <label for="exampleInputPassword1">Substream url</label>
              <input type="text" name="stream-url" class="form-control"  placeholder="Enter stream url" >
              <small  class="form-text text-muted">Enter rtsp address as instructed by your camera. Look like <code>rtsp://&lt;ip&gt;:&lt;port&gt;/path </code> </small>
            </div>
            <div class="form-group">
              <label for="inputStatus">Substream type</label>
              <select class="form-control custom-select" name="stream-ondemand" >
                <option selected disabled><small>Select One</small></option>
                <option value="1">On demand only</option>
                <option value="0">Persistent connection</option>
              </select>
              <small  class="form-text text-muted">On persistent connection, the server get data from the camera continuously. On demand, the server get data from the camera only when you click play button </small>
            </div>
            <div class="form-group">
              <div class="custom-control custom-switch">
                <input type="checkbox" class="custom-control-input" name="debug" id="substream-debug-switch-` + random + `" >
                <label class="custom-control-label" for="substream-debug-switch-` + random + `">Enable debug</label>
              </div>
              <small  class="form-text text-muted">Select this options if you want get more data about the stream </small>
            </div>
              </form>
          </div>
      </div>
      </div>`;
  return html;
}

function removeChannelDiv(element) {
  $(element).closest('.col-12').remove();
}

function logger() {
  if (!colordebug) {
    return;
  }
  let colors = {
    "0": "color:green",
    "1": "color:#66CDAA",
    "2": "color:blue",
    "3": "color:#FF1493",
    "4": "color:#40E0D0",
    "5": "color:red",
    "6": "color:red",
    "7": "color:red",
    "8": "color:red",
    "9": "color:red",
    "10": "color:red",
    "11": "color:red",
    "12": "color:red",
    "13": "color:red",
    "14": "color:red",
    "15": "color:red",
  }
  console.log('%c%s', colors[arguments[0]], new Date().toLocaleString() + " " + [].slice.call(arguments).join('|'))
}
