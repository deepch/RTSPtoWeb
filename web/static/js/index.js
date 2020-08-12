$(document).ready(() => {
  localImages();
})

$('body').on('click', '.nav-link', function() {
  $('.nav-link').removeClass('active');
  $(this).addClass('active');
})
var activeStream = null;

function showAddStream(streamName, streamUrl) {
  streamName = streamName || '';
  streamUrl = streamUrl || '';
  Swal.fire({
    title: 'Add stream',
    html: '<label>Name</label><input id="stream-name" class="swal2-input" value="' + streamName + '">' +
      '<label>Url</label><input id="stream-url" class="swal2-input" value="' + streamUrl + '">',
    focusConfirm: true,
    showCancelButton: true,
    preConfirm: () => {
      var uuid = randomUuid(),
        name = $('#stream-name').val(),
        url = $('#stream-url').val();
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
          url: url
        });
        Swal.fire(
          'Added!',
          'Your stream has been Added.',
          'success'
        )
      }
      //console.log(uuid, name, url)
    }
  })

}

function showEditStream(uuid, name, url) {

}

function deleteStream(uuid) {
  activeStream = uuid;
  console.log(activeStream);
  Swal.fire({
    title: 'Are you sure?',
    text: "Do you want delete this stream?",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#3085d6',
    cancelButtonColor: '#d33',
    confirmButtonText: 'Yes, delete it!'
  }).then((result) => {
    if (result.value) {
      goRequest('delete', activeStream)

    }
  })
}



function goRequest(method, uuid, data) {
  data = data || null;
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
  console.log(path);
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
    data: JSON.stringify(data),
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
    ajaxParam.dataType = 'json';
  }
  console.log(ajaxParam);
  $.ajax(ajaxParam);
}

function goRequestHandle(method, response, uuid) {
  console.log(method, response);
  switch (method) {
    case 'add':

      break;
    case 'edit':

      break;
    case 'delete':
      $('#' + uuid).remove();
      Swal.fire(
        'Deleted!',
        'Your stream has been deleted.',
        'success'
      )
      break;
    case 'reload':

      break;
    case 'info':

      break;
    case 'streams':

      break;
    default:

  }

}

function startPlay(type, uuid) {
  activeStream = uuid;
  $("#videoPlayer")[0].src = '';
  $("#videoPlayer")[0].load();

  switch (type) {
    case 'hls':
      playHls(uuid);
      break;
    default:
      Swal.fire(
        'Sorry',
        'This option is still under development',
        'question'
      )
      return;
  }
  $('#player-wrapper').removeClass('d-none');
}
var hls = null;
if (Hls.isSupported()) {
  hls = new Hls();
}

function playHls(uuid) {
  var url = '/stream/' + uuid + '/hls/live/index.m3u8';
  if ($("#videoPlayer")[0].canPlayType('application/vnd.apple.mpegurl')) {
    $("#videoPlayer")[0].src = url;
    $("#videoPlayer")[0].load();
  } else {
    if (hls != -null) {
      hls.loadSource(url);
      hls.attachMedia($("#videoPlayer")[0]);
    } else {
      Swal.fire({
        icon: 'error',
        title: 'Oops...',
        text: 'Your browser don`t support hls '
      })
    }
    //console.log($("#videoPlayer")[0].canPlayType('application/vnd.apple.mpegurl'))

  }
}

$("#videoPlayer")[0].addEventListener('loadeddata', function() {
  console.log('loadeddata');
  makePic();
});

function makePic() {
  ratio = $("#videoPlayer")[0].videoWidth / $("#videoPlayer")[0].videoHeight;
  w = 200;
  h = parseInt(w / ratio, 10);
  $('#canvas')[0].width = w;
  $('#canvas')[0].height = h;
  $('#canvas')[0].getContext('2d').fillRect(0, 0, w, h);
  $('#canvas')[0].getContext('2d').drawImage($("#videoPlayer")[0], 0, 0, w, h);
  var imageData = $('#canvas')[0].toDataURL();
  var images = localStorage.getItem('images');
  console.log(images);
  if (images != null) {
    images = JSON.parse(images);
  } else {
    images = {};
  }
  images[activeStream] = imageData;
  localStorage.setItem('images', JSON.stringify(images));
  $('#' + activeStream).find('.stream-img').attr('src', imageData);
  //console.log($('#canvas')[0].toDataURL());
}

function localImages() {
  var images = localStorage.getItem('images');
  if (images != null) {
    images = JSON.parse(images);
    $.each(images, function(k, v) {
      $('#' + k).find('.stream-img').attr('src', v);
    });
  }
}

function randomUuid() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

function validURL(str) {
  var pattern = new RegExp('^(rtsp?:\\/\\/)?' + // protocol
    '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
    '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
    '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
    '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
    '(\\#[-a-z\\d_]*)?$', 'i'); // fragment locator
  return !!pattern.test(str);
}