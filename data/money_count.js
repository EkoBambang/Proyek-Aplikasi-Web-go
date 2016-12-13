function numberWithCommas(x) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function convertfromcomma(x) {
    return x.toString().replace(/,/g, '');
}


function hitung(total_semua_belanja,uang_ku,kembalian,uang_diberikan){

var a = document.getElementById(total_semua_belanja).innerHTML;

var b = document.getElementById(uang_ku).value;

if (b == "") {

alert("harap mengisi kolom diatas!");
return;

}


var c = parseInt(b) - parseFloat(convertfromcomma(a));


if (c < 0) {

alert("uang anda tidak cukup!!");
return;

}
document.getElementById(uang_ku).value = "";

if (c == 0){
document.getElementById(kembalian).innerHTML = "0";

} else {

document.getElementById(kembalian).innerHTML = numberWithCommas(c);
}

document.getElementById(uang_diberikan).innerHTML = numberWithCommas(b);




}