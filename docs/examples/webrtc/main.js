document.addEventListener('DOMContentLoaded', function () {
  function startPlay (videoEl, url) {
    const webrtc = new RTCPeerConnection({
      iceServers: [{
        urls: ['stun:stun.l.google.com:19302']
      }],
      sdpSemantics: 'unified-plan'
    })
    webrtc.ontrack = function (event) {
      console.log(event.streams.length + ' track is delivered')
      videoEl.srcObject = event.streams[0]
      videoEl.play()
    }
    webrtc.addTransceiver('video', { direction: 'sendrecv' })
    webrtc.onnegotiationneeded = async function handleNegotiationNeeded () {
      const offer = await webrtc.createOffer()

      await webrtc.setLocalDescription(offer)

      fetch(url, {
        method: 'POST',
        body: new URLSearchParams({ data: btoa(webrtc.localDescription.sdp) })
      })
        .then(response => response.text())
        .then(data => {
          try {
            webrtc.setRemoteDescription(
              new RTCSessionDescription({ type: 'answer', sdp: atob(data) })
            )
          } catch (e) {
            console.warn(e)
          }
        })
    }

    const webrtcSendChannel = webrtc.createDataChannel('rtsptowebSendChannel')
    webrtcSendChannel.onopen = (event) => {
      console.log(`${webrtcSendChannel.label} has opened`)
      webrtcSendChannel.send('ping')
    }
    webrtcSendChannel.onclose = (_event) => {
      console.log(`${webrtcSendChannel.label} has closed`)
      startPlay(videoEl, url)
    }
    webrtcSendChannel.onmessage = event => console.log(event.data)
  }

  const videoEl = document.querySelector('#webrtc-video')
  const webrtcUrl = document.querySelector('#webrtc-url').value

  startPlay(videoEl, webrtcUrl)
})
