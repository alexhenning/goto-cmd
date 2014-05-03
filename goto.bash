function goto() {
    if [ "$1" == "-complete" ]
    then
        goto-cmd $1 $2
    elif [ "$1" == "-h" ] || [ "$1" == "-help" ] || [ "$1" == "--help" ]
    then
        goto-cmd -help
    else
        cd $(goto-cmd $1)
    fi
}
function _complete_goto() {
    local word=${COMP_WORDS[COMP_CWORD]}
    COMPREPLY=($(compgen -W "`goto-cmd -complete ${word}`" -- "${word}"))
}
complete -o nospace -F _complete_goto goto
