pragma solidity ^0.7.0;
contract Hello {
    function get() public returns (string memory) {
        return "Hello World!";
    }

    uint number = 0;
    function append() public {
        number = number + 1;
    }

    function getNumber() public returns (string memory){
        if (number == 0) {
            return "0";
        }
        if (number == 1) {
            return "hello: 1";
        }
        if (number == 2) {
            return "hahaha: 2";
        }
    }
}