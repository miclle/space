'use strict';

function initPagetree() {
  var pathname = window.location.pathname.replace(/\/+$/, "");
  var current = $(".pagetree a").filter(function (index, a) {
    return pathname === a.pathname;
  });

  current.addClass("active");

  // var ul = current.parents(".collapse").addClass("show");
  // ul.prev().find('[data-toggle="collapse"]').removeClass("collapsed");
  // current.parent().next(".collapse").addClass("show");

  var $sidebar = $("#sidebar");
  var offset = $(".pagetree ul a.active").offset();
  offset && $sidebar.scrollTop(offset.top - 100);
}

$(function() {
  initPagetree();
});