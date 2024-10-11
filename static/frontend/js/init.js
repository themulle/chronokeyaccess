const datetime = document.getElementById('dateSelector');

var daytimes = [];
var timerange = "";

function setTimeRange(paramStartTime) {
  timerange = paramStartTime;
  showPinCodes(daytimes)
}

function showPinCodes(paramDaytime) {
  daytimes=[...paramDaytime];
  if (timerange != "") {
    var i=0;
    while (i < paramDaytime.length) {
      paramDaytime[i] = paramDaytime[i] + "T" + timerange;
      i++;
    }
  }

  selectedTab=$("#iontabbar").attr("selected-tab")

  $.ajax({ url: "/codestemplate", data: { daytime: paramDaytime, exactmatch: (timerange != ""), codetype: (selectedTab), }, traditional: true,  })
    .done(function (data) {
      $('#codelist').html(data);
    });
}

async function copyText() {
  var contentToCopy = "";
  $(".copyContent").each(function() {
      contentToCopy += $(this).val() + "\n";
  });
 
  contentToCopy = $("#WhatsappPrefix").val()
    + "\n\n"
    + contentToCopy
    + "\n"
    + $("#WhatsappPostfix").val()
  
    try {
      await navigator.clipboard.writeText(contentToCopy);
      alert(contentToCopy);
  } catch (err) {
      console.error("Failed to copy: ", err);
      alert("Failed to copy content to clipboard.");
  }
}


$(document).ready(function () {
  datetime.addEventListener('ionChange', function () {
    console.log('ionChange', this.value, timerange);
    showPinCodes(this.value);
  });

  $(".quick-hour").on('click', function () {
    var hour = $(this).data('hour');
    $('#hourSelect').val(hour);
    $('#hourSelect').get(0).dispatchEvent(new Event('ionChange'));
  });

  {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, '0'); // Monate von 0-11, daher +1
    const day = String(today.getDate()).padStart(2, '0');
    datetime.value = [`${year}-${month}-${day}`];
  }
  showPinCodes(datetime.value);
});