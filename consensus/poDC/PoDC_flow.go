/* 블록 제안(Proposal)

프론트 노드(Proposer)는 블록 생성을 요청하는 새로운 Tx를 수신하고 Bx를 생성하는 시점에 Qmanager 서버세 접속해서 퀀텀난수 생성을 요청
Proposer가, QRND(퀀텀난수)  = RequestQRND()
RequestQRND() ; 퀀턴난수 발생기 서버에 요청해서 받는다.

Step 0: Request 단계: 프런트노드가, QManager에게 요청한다고 해서,
Request로 명명함.


프런트 노드가 QManager로부터 정보를 받아온후, Step 1로 진입.

해킹이 _불가능한 _메시지 _전달 _방식과 _고속 _인증과 _합의가 _구현된 _하드웨어 _의존적 _호출함수들의 _모임을 _퀀텀난수 _매니징 _서버 _(QM서버)라 _명명한다 _

QM서버의 _역할에 _대해, 퀀텀난수를 _접목한 _양자확률함수기법(Quantum Random Function Mechanism) 을 _도입한다. REAP CHAIN의 _차별적인 _특징은 _QRF Mechanism에 _담겨있다.



Qmanager 통지

Qmanager는 프론트 노드의 요청을 받고 블록생성 작업과 합의를 위한 코디 및 운영위 후보군과 이들의 인식을 위한 퀀텀난수 정보를 암호화하여 이를 프론트 노드에 전송. 선정된 코디에게는 운영위 후보군 정보를 전송
 1. QManager
    Quantum 난수 발생,
    프런트노드(Proposer)에 전달 전에,, 보안을 위해서 암호화,,


ReceiveQRND ()

 selectCordinator()

 selectedCordinator  : Index  ;

Selected Cordinator 입장에서는  Qmanager로부터,, 운영위 후보군,, enode값을 전송 받는다( 즉 Validator node 목록에서 )





↓


프론트 노드의 전체 통지

프론트 노드는 Qmanager로부터 전송받은 암호화 정보로 블록 Header에 Extra data로 구성하여 전체 노드에 브로드캐스팅
pre-prepare 단계
Step 1: Pre-prepare 단계:
모든 노드에 브로드캐스 하는 단계임.
받는 것은 네트웍 브로드캐스트라, 다 받으나, 응답하는 노드는
일반노드, 코디네이터, 상임위노드만 응답함.



Validator node 가  고정된 상임위 노드 ( 13 + 1( 코디네이터노드)), 선정된 운영위 노드( 15), 운영위 후보 노드, 일반노드 이렇게 종류가 있다.
 	이 종류를 다 정의해야함.

코디 및 후보군 확인

각 노드는 전송받은 블록의 Extra data에서 자기자신의 해시를 검증하고 자신의 정보인 경우에만 복호화하여 코디와 후보군 확인
일반노드의 경우 ...


↓


운영위 등록 레이싱

운영위 후보로 확인된 운영위 후보 노드들은 선착순으로 진행되는 운영위 등록 레이싱을 진행
↓


운영위 확정

코디는 선착순 15개의 노드의 운영위를 확정. 상임위 및 운영위에게 확정결과를 통보
↓


합의 진행

상임위 및 운영위 합의 진행
↓


블록 생성

코디는 합의에 따른 블록 생성

===========================================
States:



 */
 */
 */
•	NEW ROUND: Proposer(= Front node) to send new block proposal. Validators wait for PRE-PREPARE message.
•
•	PRE-PREPARED: A validator has received PRE-PREPARE message and broadcasts PREPARE message. Then it waits for 14 + 15 of PREPARE or COMMIT messages.
•
•
•	PREPARED: A validator has received 14 + 15 of PREPARE messages and broadcasts COMMIT messages. Then it waits for 14 + 15 of COMMIT messages.



•   d-select :


•	COMMITTED: A validator has received 14 + 15 of COMMIT messages and is able to insert the proposed block into the blockchain.


    d-commited:
 */
 */
 */
•
•	FINAL COMMITTED: A new block is successfully inserted into the blockchain and the validator is ready for the next round.
•
•	ROUND CHANGE: A validator is waiting for 2F + 1 of ROUND CHANGE messages on the same proposed round number.







=============================================
 */
 */
 */
 */
 */
 */
 */
 */
 */

우리가 노드에서 블록을 생성하는 과정에서는 어떤 일들이 일어나는지 확인하기 위해 디버그할 메소드 목록은 다음과 같습니다.
•	worker.commitTransaction()
•	worker.commitNewWork()
•	ethash.mine()
•	Blockchain.insert()

•	Sequence: Sequence number of a proposal. A sequence number should be greater than all previous sequence numbers.
	Currently each proposed block height is its associated sequence number.
•	Backlog: The storage to keep future consensus messages.
•	Round state: Consensus messages of a specific sequence and round, including pre-prepare message, prepare message, and commit message.


•	Consensus proof: The commitment signatures of a block that can prove the block has gone through the consensus process.

•	Snapshot: The validator voting state from last epoch.
•	마지막 시대의 유효성 검사기 투표 상태입니다.


*/

