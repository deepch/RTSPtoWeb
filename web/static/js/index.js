$(document).ready(() => {
  localImages();
  if(localStorage.getItem('defaultPlayer')!=null){
    $('input[name=defaultPlayer]').val([localStorage.getItem('defaultPlayer')]);
  }
  //console.log(localStorage.getItem('defaultPlayer'));
})
$('input[name=defaultPlayer]').on('change',function(){
  //console.log($(this).val());
  localStorage.setItem('defaultPlayer', $(this).val());
})

var activeStream = null;

function showAddStream(streamName, streamUrl) {
  streamName = streamName || '';
  streamUrl = streamUrl || '';
  Swal.fire({
    title: 'Add stream',
    html: '<form class="text-left"> '
    +'<div class="form-group">'
    +'<label>Name</label>'
    +'<input type="text" class="form-control" id="stream-name">'
    +'<small class="form-text text-muted"></small>'
    +'</div>'
    +'<div class="form-group">'
    +'  <label>URL</label>'
    +'  <input type="text" class="form-control" id="stream-url">'
    +'  </div>'
+'<div class="form-group form-check">'
    +'<input type="checkbox" class="form-check-input" id="stream-ondemand">'
    +'<label class="form-check-label">ondemand</label>'
  +'</div>'
+'</form>',
    focusConfirm: true,
    showCancelButton: true,
    preConfirm: () => {
      var uuid = randomUuid(),
        name = $('#stream-name').val(),
        url = $('#stream-url').val(),
        ondemand=$('#stream-ondemand').val();
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
          ondemand:ondemand
        });


      }
      //console.log(uuid, name, url)
    }
  })

}

function showEditStream(uuid) {
 console.log(streams[uuid]);
}

function deleteStream(uuid) {
  activeStream = uuid;
  //console.log(activeStream);
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
      goRequest('delete', uuid)

    }
  })
}

function renewStreamlist(){
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
    dataType:'json',
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
    //ajaxParam.dataType = 'json';
  }
  $.ajax(ajaxParam);
}

function goRequestHandle(method, response, uuid) {
  switch (method) {
    case 'add':

      if(response.status==1){
        renewStreamlist();
        Swal.fire(
          'Added!',
          'Your stream has been Added.',
          'success'
        );
      }else{
        Swal.fire({
          icon: 'error',
          title: 'Oops...',
          text: 'Same mistake issset',
        })
      }

      break;
    case 'edit':

      break;
    case 'delete':

      if(response.status==1){
        $('#' + uuid).remove();
        delete(streams[uuid]);
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
      if(response.status==1){
        streams=response.payload;
        if(Object.keys(streams).length>0){
          $.each(streams,function(uuid,param){
            if($('#'+uuid).length==0){
              $('.streams').append(streamHtmlTemplate(uuid,param.name));
            }
          })
        }
      }
      //console.log(streams);
      break;
    default:

  }

}

if($("#videoPlayer").length>0){
  $("#videoPlayer")[0].addEventListener('loadeddata', function() {
    console.log('loadeddata');
    makePic();
  });
}


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
  if (images != null) {
    images = JSON.parse(images);
  } else {
    images = {};
  }

  images[rtspPlayer.uuid] = imageData;
  localStorage.setItem('images', JSON.stringify(images));
  $('#' + rtspPlayer.uuid).find('.stream-img').attr('src', imageData);
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
function clearLocalImg(){
  localStorage.setItem('images','{}');
}

function streamHtmlTemplate(uuid,name){
  return '<div class="item" id="'+uuid+'">'
    +'<div class="stream">'
      +'<div class="thumbs" onclick="rtspPlayer.livePlayer(0, \''+uuid+'\')">'
      +'<img src="../static/img/noimage.svg" alt="" class="stream-img">'
      +'</div>'
      +'<div class="text">'
        +'<h5>'+name+'</h5>'
        +'<p>property</p>'
        +'<div class="input-group-prepend dropleft text-muted">'
          +'<a class="btn" data-toggle="dropdown" >'
            +'<i class="fas fa-ellipsis-v"></i>'
          +'</a>'
          +'<div class="dropdown-menu">'
            +'<a class="dropdown-item" onclick="rtspPlayer.livePlayer(\'hls\', \''+uuid+'\')" href="#">Play HLS</a>'
            +'<a class="dropdown-item" onclick="rtspPlayer.livePlayer(\'mse\', \''+uuid+'\')" href="#">Play MSE</a>'
            +'<a class="dropdown-item" onclick="rtspPlayer.livePlayer(\'webrtc\', \''+uuid+'\')" href="#">Play WebRTC</a>'
            +'<div class="dropdown-divider"></div>'
            +'<a class="dropdown-item" onclick="showEditStream(\''+uuid+'\')" href="#">Edit</a>'
            +'<a class="dropdown-item" onclick="deleteStream(\''+uuid+'\')" href="#">Delete</a>'
          +'</div>'
        +'</div>'
      +'</div>'
    +'</div>'
  +'</div>';
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
