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
