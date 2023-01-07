build:
	go build -o yara_cleaner .

clear:
	rm -f yara_cleaner; rm -rf io/cleaned/

run:
	./yara_cleaner -output io/cleaned -stripMeta -stripTags -recursive io/raw/

test:
	./yara_cleaner -output io/cleaned -stripMeta -stripTags -recursive io/test/

scan:
	yara outputDir/crypto/crypto_signatures.yar outputDir/crypto/crypto_signatures.yar > scan.txt