# normalize

    go get github.com/chbrown/normalize


## Alternatives

    find . -type f -print0 | xargs -0 sed -i '' 's/[[:space:]]*$//'

    git ls-tree -r master --name-only -z | xargs -0 sed -i '' 's/[[:space:]]*$//'


## Development

To install locally:

    export GOBIN=${GOPATH-~/go}/bin
    go install


## License

Copyright 2017–2018 Christopher Brown.
[MIT Licensed](https://chbrown.github.io/licenses/MIT/#2017-2018).
