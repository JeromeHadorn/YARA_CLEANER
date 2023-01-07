build:
	go build -o yara_cleaner .

clear:
	rm -f yara_cleaner; rm -rf output

run:
	./yara_cleaner -encode mypassword123456 -stripMeta -stripTags -recursive -output  output/ data/

scan:
	yara output/crypto/crypto_signatures.yar output/crypto/crypto_signatures.yar > scan.txt