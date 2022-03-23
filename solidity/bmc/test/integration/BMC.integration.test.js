const BMCManagement = artifacts.require('BMCManagement');
const BMCPeriphery = artifacts.require('BMCPeriphery');
const MockBMCPeriphery = artifacts.require('MockBMCPeriphery');
const MockBMCManagement = artifacts.require('MockBMCManagement');
const MockBSH = artifacts.require('MockBSH');
const { assert } = require('chai');
const truffleAssert = require('truffle-assertions');
const URLSafeBase64 = require('urlsafe-base64');
const rlp = require('rlp');

const { deployProxy } = require('@openzeppelin/truffle-upgrades');

contract('BMC tests', (accounts) => {

    describe('Handle relay message tests', () => {
        let bmcManagement, bmcPeriphery, bmv, bsh; 
        let network = '0x7.icon';
        let link = "btp://0x7.icon/cxfe6b306c41bf7cd880dafe46a952fb4d1764d49b"
        let height = 0;
        let offset = 0;
        let lastHeight = 0;
        let blockInterval = 3000;
        let maxAggregation = 5;
        let delayLimit = 3;
        let relays;
    
        beforeEach(async () => {
            bmcManagement = await deployProxy(BMCManagement);
            console.log('bmcManagement:', bmcManagement.address);
            bmcPeriphery = await deployProxy(MockBMCPeriphery, ['0x97.bsc', bmcManagement.address]);
            
            console.log('bmcPeriphery:', bmcPeriphery.address);
            await bmcManagement.setBMCPeriphery(bmcPeriphery.address);
            
            bsh = await MockBSH.new();
            console.log('bsh:', bsh.address);
            await bmcManagement.addService('TokenBSH', bsh.address);
            await bmcManagement.addLink(link);
            relays = [accounts[0]];
            await bmcManagement.addRelay(link, relays);
            await bmcManagement.setLink(link, blockInterval, maxAggregation, delayLimit);
           
        });

        it('Scenario 1: Revert if relay is invalid', async() => {
            const btpMsg = rlp.encode([
                'btp://0x03.icon/cx10c8c08724e7a95c84829c07239ae2b839a262a3',
                'btp://0x27.pra/' + bmcPeriphery.address,
                'Token',
                '0x01', // rlp encode of signed int
                'message'
            ]);
           
            let _relayMsg = "-QE--QE7uQE4-QE1AbkBLfkBKvkBJ7g5YnRwOi8vMHg2MS5ic2MvMHhGODU5ODUyODgwZDczMEY3RjU3NTNFMDY2MGUzMUNjOTkwN2Y1QjlFGLjp-Oe4OWJ0cDovLzB4Ny5pY29uL2N4NzgwNGQyMzc2YjhiYzg2MWUxZjI1NzA5NTcyNTQ4NWVkZTY4ZTY3Nbg5YnRwOi8vMHg2MS5ic2MvMHhGODU5ODUyODgwZDczMEY3RjU3NTNFMDY2MGUzMUNjOTkwN2Y1QjlFiFRva2VuQlNIF7hl-GMAuGD4XqpoeGRmODE3ZDFlMjEzODg0MDgzMzE3YzRhOTc4Y2M1ZTIzOTM2ZThiODCqMHg1MDRjZjJkQWUwQUIzNGFDZWFBRUIyOTU0NjczNjE3ZjVCMUEyMkEwx8aDRVRIBQCDRNY8"
            
            let buff = new Buffer(_relayMsg, 'base64');
            let data = buff.toString('utf-8');
            //await bmcManagement.updateLinkRxSeq("btp://0x7.icon/cxfe6b306c41bf7cd880dafe46a952fb4d1764d49b",8)
            await bmcPeriphery.handleRelayMessage("btp://0x7.icon/cxfe6b306c41bf7cd880dafe46a952fb4d1764d49b", _relayMsg);
    
            let bmcLink = await bmcManagement.getLink(link);
            let status= await bmcPeriphery.getStatus(link);
            console.log(status)
        });
    
    });
});    

