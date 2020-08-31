var rtspPlayer={
  active:false,
  type:'live',
  hls:null,
  ws:null,
  mseSourceBuffer:null,
  mse:null,
  mseQueue:[],
  mseStreamingStarted:false,
  webrtc:null,
  webrtcSendChannel:null,
  webrtcSendChannelInterval:null,
  uuid:null,

  clearPlayer:function(){
    if(this.active){

      if(this.hls!=null){
        this.hls.destroy();
        this.hls=null;
      }
      if(this.ws!=null){
        //close WebSocket connection if opened
        this.ws.close(1000);
        this.ws=null;
      }
      if(this.webrtc!=null){
        clearInterval(this.webrtcSendChannelInterval);

        this.webrtc=null;
      }
      $('#videoPlayer')[0].src = '';
      $('#videoPlayer')[0].load();


      this.active=false;
    }
  },
  livePlayer:function(type,uuid){
    this.clearPlayer();
    this.uuid=uuid;
    this.active=true;

    $('.streams-vs-player').addClass('active-player');
    if(type==0){
      type=$('input[name=defaultPlayer]:checked').val()
    }
    switch (type) {
      case 'hls':
        this.playHls();
        break;
      case 'mse':
        this.playMse();
        break;
      case 'webrtc':
          this.playWebrtc();
          break;
      default:
        Swal.fire(
          'Sorry',
          'This option is still under development',
          'question'
        )
        return;
    }

  },
  playHls:function(){
    if(this.hls==null && Hls.isSupported()){
        this.hls = new Hls();
    }
    if ($("#videoPlayer")[0].canPlayType('application/vnd.apple.mpegurl')) {
      $("#videoPlayer")[0].src = this.streamPlayUrl('hls');
      $("#videoPlayer")[0].load();
    } else {
      if (this.hls != null) {
        this.hls.loadSource(this.streamPlayUrl('hls'));
        this.hls.attachMedia($("#videoPlayer")[0]);
      } else {
        Swal.fire({
          icon: 'error',
          title: 'Oops...',
          text: 'Your browser don`t support hls '
        })
      }
    }
  },
  playWebrtc:function(){
    var _this=this;
    this.webrtc=new RTCPeerConnection({
      iceServers: [{
        urls: ["stun:stun.l.google.com:19302"]
      }]
    });
    this.webrtc.onnegotiationneeded = this.handleNegotiationNeeded;
    this.webrtc.ontrack = function(event) {
      console.log(event.streams.length + ' track is delivered');
      $("#videoPlayer")[0].srcObject = event.streams[0];
      $("#videoPlayer")[0].play();
    }
    this.webrtc.addTransceiver('video', {
      'direction': 'sendrecv'
    });
    this.webrtcSendChannel = this.webrtc.createDataChannel('foo');
    this.webrtcSendChannel.onclose = () => console.log('sendChannel has closed');
    this.webrtcSendChannel.onopen = () => {
      console.log('sendChannel has opened');
      this.webrtcSendChannel.send('ping');
      this.webrtcSendChannelInterval =  setInterval(() => {
        this.webrtcSendChannel.send('ping');
      }, 1000)
    }

    this.webrtcSendChannel.onmessage = e => console.log(e.data);
  },
  handleNegotiationNeeded: async function(){
    var _this=rtspPlayer;

    offer = await _this.webrtc.createOffer();
    await _this.webrtc.setLocalDescription(offer);
    $.post(_this.streamPlayUrl('webrtc'), {
      data: btoa(_this.webrtc.localDescription.sdp)
    }, function(data) {
      //console.log(data)
      try {

        _this.webrtc.setRemoteDescription(new RTCSessionDescription({
          type: 'answer',
          sdp: atob(data)
        }))



      } catch (e) {
        console.warn(e);
      }

    });
  },
  playMse:function(){
    //console.log(this.streamPlayUrl('mse'));
    var _this=this;
    this.mse = new MediaSource();
    $("#videoPlayer")[0].src=window.URL.createObjectURL(this.mse);
    this.mse.addEventListener('sourceopen', function(){
      _this.ws=new WebSocket(_this.streamPlayUrl('mse'));
      _this.ws.binaryType = "arraybuffer";
      _this.ws.onopen = function(event) {
        console.log('Connect to ws');
      }

      _this.ws.onmessage = function(event) {
        var data = new Uint8Array(event.data);
        if (data[0] == 9) {
          decoded_arr=data.slice(1);
          if (window.TextDecoder) {
            mimeCodec = new TextDecoder("utf-8").decode(decoded_arr);
          } else {
            mimeCodec = Utf8ArrayToStr(decoded_arr);
          }
          console.log(mimeCodec);
          _this.mseSourceBuffer = _this.mse.addSourceBuffer('video/mp4; codecs="' + mimeCodec + '"');
          _this.mseSourceBuffer.mode = "segments"
          _this.mseSourceBuffer.addEventListener("updateend", _this.pushPacket);

        } else {
          _this.readPacket(event.data);
        }
      };
    }, false);

  },
  readPacket:function(packet){
    if (!this.mseStreamingStarted) {
      this.mseSourceBuffer.appendBuffer(packet);
      this.mseStreamingStarted = true;
      return;
    }
    this.mseQueue.push(packet);

    if (!this.mseSourceBuffer.updating) {
      this.pushPacket();
    }
  },
  pushPacket:function(){
    var _this=rtspPlayer;
    if (!_this.mseSourceBuffer.updating) {
      if (_this.mseQueue.length > 0) {
        packet = _this.mseQueue.shift();
        var view = new Uint8Array(packet);
        _this.mseSourceBuffer.appendBuffer(packet);
      } else {
        _this.mseStreamingStarted = false;
      }
    }
  },
  streamPlayUrl:function(type){
    switch (type) {
      case 'hls':
        return '/stream/' + this.uuid + '/hls/live/index.m3u8';
        break;
      case 'mse':
        var potocol = 'ws';
        if (location.protocol == 'https:') {
          potocol = 'wss';
        }
        return potocol+'://'+location.host+'/stream/' + this.uuid +'/mse?uuid='+this.uuid;
        //return 'ws://sr4.ipeye.ru/ws/mp4/live?name=d4ee855e40874ef7b7149357a42f18f0';
        break;
      case 'webrtc':
        return  "/stream/"+this.uuid+"/webrtc?uuid=" + this.uuid;
        break;
      default:
        return '';
    }
  }

}
