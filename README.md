# To run

docker-compose up --build
examples:

1. GetFileContents(pathInRepo) -> fileContents
   curl http://localhost:3334/file-contents/README.md

2. CheckoutRef(gitRef) -> [no reply]
   curl http://localhost:3334/checkout-ref/master

3. HashFiles(pathInRepo, ...) -> hash
   curl -X POST http://localhost:3333/hash-files/echo/post/json -H "Content-Type: application/json" -d '{"paths": ["README.md", ".gitignore"]}'
