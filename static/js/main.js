/**
 * Created by Administrator on 2017/5/8.
 */

$(function(){

    var evt1 = 'ontouchmove' in window ? 'touchmove' : 'mousemove';
    window.addEventListener(evt1, function(e){
        e.preventDefault();

    });

    /*首页*/
    var client_height=document.documentElement.clientHeight;
    var client_width=document.documentElement.clientWidth;
    $("#index,#map,#page3,#page4,#page5,#page6,#page7").height(client_height);
    $("#index,#map,#page3,#page4,#page5,#page6,#page7").width(client_width);

    $(".fluid").click(function(){
        window.location.href="login.html"
    })

    /*page2*/
    $("#map .button").click(function(){
        window.location.href="stake.html"
    })

    /*page3*/
    /*$("#page3").click(function(){
        window.location.href="page4.html"
    })*/
    $(".weather1").click(function(){
        $(".number-money").text(100);
        $(".number-weather").text(10);
        $(".number-cow").text(19);
        $(".number-date").text("10天");
    })
    $(".weather2").click(function(){
        $(".number-money").text(90);
        $(".number-weather").text(8);
        $(".number-cow").text(12);
        $(".number-date").text("8天");
    })
    $(".weather3").click(function(){
        $(".number-money").text(60);
        $(".number-weather").text(5);
        $(".number-cow").text(15);
        $(".number-date").text("18天");
    })
    $(".weather4").click(function(){
        $(".number-money").text(200);
        $(".number-weather").text(20);
        $(".number-cow").text(30);
        $(".number-date").text("13天");
    })

    /*page4*/
    $("#page4").click(function(){
        window.location.href="lose.html"
    })

    /*page5*/
    $("#page5").click(function(){
        window.location.href="rank.html"
    })

    /*page6*/
    $("#page6").click(function(){
        window.location.href="end.html"
    })

})



