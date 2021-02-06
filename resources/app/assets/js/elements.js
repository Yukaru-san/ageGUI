// Stores all the elements containing an ID

// File Area
var uploader = $('.plupload.html5')[0]

// Table
var fileTable = document.getElementById('filePreview');

// Button
var continueBtn = document.getElementById("continueBtn");
var disableBtn = document.getElementById("disableBtn");

// Pages
var uploadPage = document.getElementById("uploadPage");
var handlePage = document.getElementById("handlePage");

// Form
var cryptChoice = document.getElementById("cryptChoice");
var zipFilesParent = document.getElementById("zipFilesParent");
var zipFiles = document.getElementById("zipFiles");
var useArmorParent = document.getElementById("useArmorParent");
var useArmor = document.getElementById("useArmor");
var key = document.getElementById("keyEntry");
var passwordInputToggle = document.getElementById("passwordInputToggle");
var keyInputToggle = document.getElementById("keyInputToggle");
var keyEntry = document.getElementById("keyEntry");
var passwordEntry = document.getElementById("passwordEntry");
var keyEntryFile = document.getElementById("keyEntryFile");
var fileInputToggle = document.getElementById("fileInputToggle");
var output = document.getElementById("outputPath");

// Overlays
var loadingOverlay = document.getElementById("loadingOverlay");
var keyReturnOverlay = document.getElementById("keyReturnOverlay");
var keyReturnOutputPath = document.getElementById("keyReturnOutputPath");
var alertOverlay = document.getElementById("alertOverlay");
var alertDiv = document.getElementById("alertDiv");
var alertTextBold = document.getElementById("alertTextBold");
var alertText = document.getElementById("alertText");
