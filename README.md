# BM
Book management

make create-book TITLE="The Great Gatsby" AUTHOR="F. Scott Fitzgerald" GENRE="Classic" PUBLISHED_DATE="2023-01-01"
make update-book ID=1 TITLE="Updated Title" AUTHOR="Updated Author" GENRE="Updated Genre"
make delete-book ID=1


make get-collections IDS=1,2,3 ORDER_BY=name DESC=true PAGE=1 PAGE_SIZE=10
make create-collection NAME="New Collection" DESCRIPTION="A new collection"
make update-collection ID=1 NAME="Updated Collection" DESCRIPTION="An updated collection"
make delete-collection ID=1


make create-books-collection CID=1 BOOK_IDS=2,3,4
make delete-books-collection CID=1 BOOK_IDS=2,3,4
