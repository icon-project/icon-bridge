packages=("bmc" "bsh/nep141" "bsh/bts")

for index in ${!packages[@]};
do 
    chmod +x "$PWD/${packages[$index]}/compile.sh"
    "$PWD/${packages[$index]}/compile.sh"
done