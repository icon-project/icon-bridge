remove:
	docker-compose down -v --remove-orphans	
	docker rmi -f icon-bsc_btp
	docker rm -f javascore-dist
	docker rmi -f iconbridge
	docker rmi -f btp/javascore
	rm -rf local/artifacts local/log
	rm -rf ../../../build/*
	rm -rf ./data/bsc/node1
	
clean:
	docker-compose down -v --remove-orphans
	rm -rf local/artifacts

build:
	make -C ../../../ dist-javascore dist-sol iconbridge-image
	mkdir -p ./data/bsc/node1
	docker-compose build

run-bsc:
	docker-compose up --remove-orphans -d binancesmartchain

run:
	docker-compose up --remove-orphans -d

build-testnet:
	make -C ../../../ dist-javascore dist-sol iconbridge-image
	docker-compose -f docker-compose.testnet.yml build

run-testnet:
	docker-compose -f docker-compose.testnet.yml up --remove-orphans -d

solidity:
	make -C ../../../ dist-sol iconbridge-image
	docker-compose up --remove-orphans -d binancesmartchain

alice_transfer_eth:
	docker-compose exec btp sh /iconbridge/bin/transfer_eth.sh

alice_transfer_icx:
	docker-compose exec btp sh /iconbridge/bin/transfer_icx.sh

alice_transfer_bnb:
	docker-compose exec btp sh /iconbridge/bin/transfer_bnb.sh

bob_transfer_eth:
	docker-compose exec btp sh /iconbridge/bin/transfer_eth_bsc.sh

bob_transfer_bnb:
	docker-compose exec btp sh /iconbridge/bin/transfer_bnb_bsc.sh

bob_transfer_icx:
	docker-compose exec btp sh /iconbridge/bin/transfer_icx_bsc.sh


testnet_alice_transfer_eth:
	docker-compose -f docker-compose.testnet.yml exec btp_testnet sh /iconbridge/bin/transfer_eth.sh

testnet_alice_transfer_icx:
	docker-compose -f docker-compose.testnet.yml exec btp_testnet sh /iconbridge/bin/transfer_icx.sh

testnet_alice_transfer_bnb:
	docker-compose -f docker-compose.testnet.yml exec btp_testnet sh /iconbridge/bin/transfer_bnb.sh

testnet_bob_transfer_eth:
	docker-compose -f docker-compose.testnet.yml exec btp_testnet sh /iconbridge/bin/transfer_eth_bsc.sh

testnet_bob_transfer_bnb:
	docker-compose -f docker-compose.testnet.yml exec btp_testnet sh /iconbridge/bin/transfer_bnb_bsc.sh

testnet_bob_transfer_icx:
	docker-compose -f docker-compose.testnet.yml exec btp_testnet sh /iconbridge/bin/transfer_icx_bsc.sh