var n = new Date();
 y = n.getFullYear();
 m = n.getMonth() + 1;
 d = n.getDate();

n = y + "-" + m + "-" + d;
document.getElementById("selDate").value = n;