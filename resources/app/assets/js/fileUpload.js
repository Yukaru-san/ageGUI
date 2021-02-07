var addedFiles = [];
var filePaths = [];

// Create the uploader
var uploader = new plupload.Uploader({
  runtimes: 'html5',
  drop_element: 'drop-target',
  browse_button: 'drop-target',
  max_file_size: '0',
  upload: "upload.php"
});

// Bind the upload functions on the area
uploader.bind('Init', function(up, params) {
  if (uploader.features.dragdrop) {
    
    var dropArea = document.getElementById("drop-target");

    dropArea.ondragover = function(event) {
        event.dataTransfer.dropEffect = "copy";
    };

    dropArea.ondragenter = function() {
        this.className = "dragover";
    };

    dropArea.ondragleave = function() {
        this.className = "";
    };

    dropArea.ondrop = function(x) {
        this.className = "";
    };
  }
});

// Handles files being added
function filesAdded(up, files) {
  console.log(files[0]);
  for (var i in files) {
    addedFiles.push(files[i]);
  }
  
  showFiles();
}

// Shows all files within addedFiles array
function showFiles() {
  fileTable.innerHTML = "";
  filePaths = [];

  for (var i in addedFiles) {

    // Table Row
    var tr = document.createElement("tr");

    // Inner
    var th = document.createElement("th");
    th.innerHTML = i;
    th.scope = "row";

    var td_name = document.createElement("td");
    td_name.className = "d-inline-block col-7 nameCol";
    $(td_name).attr("data-toggle", "tooltip");

    // Add file's paths
    var filePath = "";
    if (basePath == "")
      filePath = addedFiles[i].relativePath;
    else 
      filePath = path.relative(basePath, addedFiles[i].relativePath);

    $(td_name).attr("title", filePath);
    filePaths.push(filePath);

    // file's name
    var split = filePath.split("\\");
    var fileName = split[split.length-1];
    var maxNameSize = 37 - Math.round((2 * i / 100));
    if (fileName.length > maxNameSize)
      td_name.innerHTML = fileName.substr(0, maxNameSize) + "...";
    else 
      td_name.innerHTML = fileName;

    // file size
    var td_size = document.createElement("td");
    td_size.className = "d-inline-block col-2";
    td_size.innerHTML = plupload.formatSize(addedFiles[i].size);

    // Remove file
    var td_remove = document.createElement("td");
    td_remove.className = "d-inline-block col-3 textCenter";

    // Remove file btn action
    var td_remove_btn = document.createElement("button");
    td_remove_btn.className = "hiddenButton";
    td_remove_btn.innerHTML = "Delete";
    td_remove_btn.connectedFileID = i;
    td_remove_btn.addEventListener("click", function(event) {
      addedFiles.splice(this.connectedFileID, 1);
      console.log("Removed "+filePaths[this.connectedFileID]);
      showFiles()
    });

    td_remove.appendChild(td_remove_btn);
    
    // Add to table row
    tr.appendChild(th);
    tr.appendChild(td_name);
    tr.appendChild(td_size);
    tr.appendChild(td_remove);

    fileTable.appendChild(tr);
  }

  // Prepare Tooltips
  $(document).ready(function(){
    $('[data-toggle="tooltip"]').tooltip();   
  });

  // Continue Btn handling
  if (addedFiles.length == 0) {
    continueBtn.classList.toggle("disabled");
    continueBtn.disabled = true;
  }
  else if (continueBtn.classList.contains("disabled")) {
    continueBtn.classList.toggle("disabled");
    continueBtn.disabled = false;
  }

  // Delete All Btn handling
  if (addedFiles.length >= 10) 
    disableBtn.style.display = "block";
  else 
    disableBtn.style.display = "none";
}

// Bind fileAdded handler and initiate
uploader.bind('FilesAdded', filesAdded);
uploader.init();

// Used from the HTML-File, toggles views
function togglePageVisibility() {

  // Screen 1 -> Screen 2
  if (uploadPage.style.display === "none") {
    uploadPage.style.display = "block";
    uploader.style.display = "block";
    handlePage.style.display = "none";
  // Screen 2 -> Screen 1
  } else {
    uploadPage.style.display = "none";
    uploader.style.display = "none";
    handlePage.style.display = "block";
  }
}

// Toggles the visibility of the Zip and Armor Selection
function toggleEncryptSelect() {
  // Make visible
  if (zipFilesParent.style.display === "none") {
    zipFilesParent.style.display = "block";
    useArmorParent.style.display = "block";
    keyEntry.setAttribute("data-original-title", "Generates Key Pair if empty");
    keyEntry.placeholder = "age1aerfzuo7j8907j8k90defrwnzumiofznj4rf...";
  // Hide
  } else {
    zipFilesParent.style.display = "none";
    useArmorParent.style.display = "none";
    keyEntry.setAttribute("data-original-title", "");
    keyEntry.placeholder = "AGE-SECRET-KEY-1Q36W52MMHRMK0CXF7A...";
  }
}

// Swaps between a password entry field and a key one's
function swapKeyPasswordSelect(caller) {
  switch (caller) {
    case "keyInputToggle": 
      keyEntry.style.display = "block";
      passwordEntry.style.display = "none";
      keyEntryFile.style.display = "none";
      keyInputToggle.classList.remove("darkerText");
      passwordInputToggle.classList.add("darkerText");
       fileInputToggle.classList.add("darkerText");
      break;
    case "passwordInputToggle": 
      passwordEntry.style.display = "block";
      keyEntry.style.display = "none";
      keyEntryFile.style.display = "none";
      passwordInputToggle.classList.remove("darkerText");
      keyInputToggle.classList.add("darkerText");
      fileInputToggle.classList.add("darkerText");
      break;
    case "fileInputToggle": 
      keyEntry.style.display = "none";
      passwordEntry.style.display = "none";
      keyEntryFile.style.display = "block";
      fileInputToggle.classList.remove("darkerText");
      passwordInputToggle.classList.add("darkerText");
      keyInputToggle.classList.add("darkerText");
      break;
  }
}

// Removes all files given on the front page and refreshes
function removeAllFiles() {
  addedFiles = [];
  filePaths = [];
  showFiles();
}

// Sends the form to GO to handle encryption / decryption
function sendForm() {
  
  // Find the path
  var encKey = "";
  if (!keyInputToggle.classList.contains("darkerText"))
    encKey = keyEntry.value;
  else if (!passwordInputToggle.classList.contains("darkerText"))
    encKey = passwordEntry.value;
  else if (!fileInputToggle.classList.contains("darkerText"))
    encKey = keyEntryFile.files[0].path;
    
  // Create and fill the JSON
  var innerJSON = {
    encrypt: cryptChoice.checked,
    zip: zipFiles.checked,
    armor: useArmor.checked,
    key: encKey,
    usePassword: passwordInputToggle.classList.contains("darkerText") ? false : true,
    output: output.value,
    paths: filePaths
  }

  var outerJSON = {
    name: "ageRequest",
    payload: innerJSON
  }

  // Send
  astilectron.sendMessage(outerJSON, function(message) {
    loadingOverlay.style.display = "none";

    // Existing key used
    if (message.payload == undefined) {
      togglePageVisibility();
    }
    // Received new Keys
    else if (message.payload.startsWith("generatedKeys")) {
      keyReturnOverlay.style.display = "block";
      var keys = message.payload.split("%");
      generatedPublicKey.value = keys[1];
      generatedPrivateKey.value = keys[2];
      keyReturnOutputPath.value = keys[3];
    }
    // Received output Path
    else if (message.payload.startsWith("outputPath")) {
      alertOverlay.style.display = "block";
      alertTextBold.innerHTML = message.payload;
      if (alertDiv.classList.contains("alert-danger")) {
        alertDiv.classList.remove("alert-danger");
        alertDiv.classList.add("alert-success");
      }
      alertTextBold.innerHTML = "Success!";
      alertText.innerHTML = "Output-Path: " +message.payload.split("%")[1];
    }
    // Error
    else {
      alertOverlay.style.display = "block";
      alertTextBold.innerHTML = message.payload;
      if (alertDiv.classList.contains("alert-success")) {
        alertDiv.classList.remove("alert-success");
        alertDiv.classList.add("alert-danger");
      }

      switch (message.payload) {
        case "inputPathError":
          alertText.innerHTML = "Can't read one or more input files. Did you delete them?";
          break;
        case "invalidKeyError":
          alertText.innerHTML = "Invalid Key provided. Check your input.";
          break;  
        case "invalidPasswordError":
          alertText.innerHTML = "Invalid Password provided. Check your input.";
          break;
        case "writeError":
          alertText.innerHTML = "Please check your output path.";
          break;
        case "zipError":
          alertText.innerHTML = "Please check your output path.";
          break;
        default:
          alertTextBold.innerHTML = "";
          alertText.innerHTML = message.payload;
          break;
      }
    }

  });

  // Block View as long as it is busy
  loadingOverlay.style.display = "block";
}