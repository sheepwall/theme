Set XResources settings in a file, from different input types and run the new file.

### Install

Use ``go install``. 
(Optional) Rename the file to avoid conflicts or to your taste.
(Optional) Add the file to your path.

### Usage

    theme help              Returns short help and usage
    theme set filename      Set current colors to colors in "filename"
    theme save filename     Save current colors to file "filename"

### Filetypes
Currently, only XResources files (content like: ``*.color0: #ABCDEF``) and "pywal" JSON files (found in [this repo](https://github.com/dylanaraps/pywal)) are supported. I will add more when I feel the need to.
