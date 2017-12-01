build:
	rm -fr public
	gean build

publish:
	make build
	cp -i CNAME public
	git add .
	git commit -n -m update
	git push origin master
	git subtree push --prefix=public origin gh-pages