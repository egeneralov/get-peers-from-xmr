on:
  push:
    tags:
    - 'v*'

name: Release

jobs:
  build:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@master
      
      - name: build linux
        run: GOOS=linux go build -v -o get-peers-from-xmr-linux .
      
      - name: build darwin
        run: GOOS=darwin go build -v -o get-peers-from-xmr-darwin .
      
      - name: Create Release
        id: create_release
        uses: actions/create-release@v1.0.0
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: Release ${{ github.ref }}
          draft: false
          prerelease: false
      
      - name: save linux artifact in Actions
        uses: actions/upload-artifact@v1.0.0
        with:
          name: linux
          path: get-peers-from-xmr-linux
      
      - name: save darwin artifact in Actions
        uses: actions/upload-artifact@v1.0.0
        with:
          name: darwin
          path: get-peers-from-xmr-darwin

      - name: upload linux release asset
        uses: actions/upload-release-asset@v1.0.1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./get-peers-from-xmr-linux
          asset_name: get-peers-from-xmr-linux
          asset_content_type: application/binary
      
      - name: upload darwin release asset
        uses: actions/upload-release-asset@v1.0.1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./get-peers-from-xmr-darwin
          asset_name: get-peers-from-xmr-darwin
          asset_content_type: application/binary
