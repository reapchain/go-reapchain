Readme.txt

1. config.json을 읽어와서,
   합의에서 사용하는 criteria 값.. : 전체 노드 갯수 중에서 투표에 참여하는 51%에 해당하는 노드 갯수를 정해준다.
   bootnodes 를 여기에 할당해주면, geth가 노드 기동시 읽어서, 설정해주는데, bootstrap node log 보면,
   제대로 읽어오는지 확인 가능함.

2. 립체인 관련 설정값을 여기에 지정한다.

3. 향후 토큰 관련 설정 값도 여기에 지정한다.


{
  "local": {
    "consensus": {
      "criteria": 3
    },
    "token": {

    },
    "bootnodes": [
      "enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@127.0.0.1:30301"
    ]
    "senatornodes":[
     "0x2bb9d58ab564c62b9c64f29d326a4d841871f229"
     "0xca81cd471a007808bc15185357a96e9f09dcb394"
     "0x3f1434e81472a3bdcffaf1804e2ade2ae7c9f9dc"
     "0xc6a0a59f3188fc28e6aac94ba48cbba257732acb"
     "0x0ab49e6d39c544cb3ec89633504c08c93fe56802"
     "0x9a31473d87344330179c6f58c1592ea5bd20e84b"
     "0x572be5121433bfe4b75219306ab32b93ea0a95c2"
     "0xdff4dc7d2a16af8e046685d24276f70ed7dac466"
     "0x7bac404f83c40d61a6a272018b2a48d8e1ac7af7"
     "0xa670532a929485b810da7a800781381da5267d2a"
     "0xb259b5da8b0908011a8ff74078c61067e4820b4a"
     "0xb29ac8821b52f63d2d43f9f1612a73bd82d6466d"
     "0x0a4108b94cfac69bf0665f90534c35e637a4b6db"
     "0x39fd383d85dac1850b4b1986d292de51df5ffe8a"
    ]

  },
  "development": {
    "consensus": {
      "criteria": 3
    },
    "token": {

    },
    "bootnodes": [
      "enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@192.168.0.100:30301"
    ]
  },
  "production": {
    "consensus": {
      "criteria": 3
    },
    "token": {

    },
    "bootnodes": [
      "enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16:dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@192.168.0.100:30301"
    ]
  }
}
