const placeForGroupie = document.getElementById('groupie')
const placeForHints = document.getElementById('hint_place')

const templateBand =  ` 
    <div class="band">
    <div class="content">
    <h4 class="artist"><a class="link" href="/artist?id=ID">NAME</a></h4>
    </div>
    <a class="link2" href="/artist?id=ID">
    <div class="black_rectangle">
    <img class="image" src="IMAGE" alt="NAME photo">
    </div>
    </a>
    </div>
    `
const templateHint = `
    <li data-id="ID"><a href="/artist?id=ID">HINT</a></li>
    `
let hints = []

searchArtists()

class Hint {
    ID 
    Text
    isSelected
    constructor(ID, Text) {
        this.ID = ID
        this.Text = Text
        this.isSelected = false
    }
}


function searchArtists(){
    const searchBar = document.getElementById('search_bar')
    displayArtists(searchBar.value)
    // console.log(artistsData)
}

function displayArtists(text){
    hints = []
    placeForGroupie.innerHTML = ''
    placeForHints.innerHTML = ''
    if (text === '') {
        artistsData.forEach(artist => {
          render(artist)
        });

    }else {
        
        artistsData.forEach(artist => {
            let isName = false
            let isMember = false
            let isFirstAlbumDate = false
            let isSince = false
            let isLocation = false
        
            if (artist.name.toLowerCase().includes(text.toLowerCase())) {
                // if (hints.length<10){
                hints.push(new Hint( artist.id,  artist.name + ' - artist' ));
                // }
                isName = true
                
            }
            if (artist.creationDate.toString().includes(text.toLowerCase())){
                // if (hints.length<10){
                hints.push(new Hint(artist.id,artist.creationDate.toString() + '/'+ artist.name+ ' - creation date' ));
                // }
                isSince = true
                
            }
            if (artist.firstAlbum.toString().includes(text.toLowerCase())){
                // if (hints.length<10){
                hints.push(new Hint( artist.id,  artist.firstAlbum.toString() + '/'+ artist.name +' - first album'));
                // }
                isFirstAlbumDate = true
                
            }
            for (let i =0 ; i<artist.members.length; i++){
                if (artist.members[i].toLowerCase().includes(text.toLowerCase())){
                    // if (hints.length<10){
                    hints.push(new Hint( artist.id, artist.members[i]  + '/'+ artist.name + ' - member' ));
                    // }
                    isMember = true
                    
                }
            }
            for (let i=0; i<artist.Geo.length; i++){
                if (artist.Geo[i].toLowerCase().includes(text.toLowerCase())){
                    // if (hints.length<10){
                    hints.push(new Hint(  artist.id,  artist.Geo[i]+ '/'+ artist.name + ' - location' ));
                    // }
                    isLocation = true
                    
                }
            }
            
            if (isMember|| isName|| isLocation || isSince || isFirstAlbumDate ){
                render(artist)
            }
            
        })
        displayHints(hints)
    }
}

function render(artist){
    const artistHtml = templateBand
    .replace(/ID/g, artist.id)
    .replace(/NAME/g, artist.name)
    .replace(/IMAGE/g, artist.image)
    placeForGroupie.insertAdjacentHTML("beforeend", artistHtml)
}

function displayHints(hints) {
    placeForHints.innerHTML = '';
    hints.sort((a, b) => customSort(a, b, searchBar.value));

    hints.forEach(hint => {

        const hintHtml = templateHint
            .replace(/HINT/g, hint.Text)
            .replace(/ID/g, hint.ID);
        placeForHints.insertAdjacentHTML('beforeend', hintHtml);
    });
    if (hints.length>0){
        hints[0].isSelected = true
        paintSelectedHint( hints[0],hints)
    }
}

const searchBar = document.getElementById('search_bar')

searchBar.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        event.preventDefault();
        const selectedHint = getSelectedHint(hints);
        if (selectedHint) {
            const hintElement = document.querySelector(`[data-id="${selectedHint.ID}"]`);
            if (hintElement) {
                hintElement.click();
            }
        }
    }
    if (event.key === 'ArrowDown') { 
        event.preventDefault(); 
        const selectedIndex = hints.findIndex(hint => hint.isSelected);
        if (selectedIndex >= 0 && selectedIndex < hints.length - 1) {
            hints[selectedIndex].isSelected = false;
            paintSelectedHint(hints[selectedIndex + 1], hints);
        }
    }

    if (event.key === 'ArrowUp') {
        event.preventDefault(); 
        const selectedIndex = hints.findIndex(hint => hint.isSelected);
        if (selectedIndex > 0) {
            hints[selectedIndex].isSelected = false;
            paintSelectedHint(hints[selectedIndex - 1], hints);
        }
    }
});


searchBar.addEventListener('blur', function() {
    setTimeout(function() {
        placeForHints.style.display = 'none';
    }, 300);
});

searchBar.addEventListener('focus', function() {
    searchArtists();
    placeForHints.style.display = 'block'
});

function getSelectedHint(hints) {
    // hints.forEach(h=> console.log(h.isSelected))
    return hints.find(hint => hint.isSelected)
}

function paintSelectedHint(selectedHint, hints) {

    hints.forEach(h => {
        const hintElement = document.querySelector(`[data-id="${h.ID}"]`);
        if (hintElement == selectedHint) {
                hintElement.style.backgroundColor = "#E0E0E0";

                
        }else{
            h.isSelected = false
            hintElement.style.backgroundColor = "transparent";
        }
    });
}
function customSort(a, b, searchText) {
    const order = ["location", "creation date", "first album", "artist", "member"];
    const orderMap = {};
    order.forEach((item, index) => {
        orderMap[item] = index;
    });

    const aIndex = orderMap[a.Text.split(" - ")[1]] || Infinity;
    const bIndex = orderMap[b.Text.split(" - ")[1]] || Infinity;

    const aStartsWithSearch = a.Text.toLowerCase().startsWith(searchText.toLowerCase());
    const bStartsWithSearch = b.Text.toLowerCase().startsWith(searchText.toLowerCase());

    if (aStartsWithSearch && !bStartsWithSearch) {
        return -1;
    } else if (!aStartsWithSearch && bStartsWithSearch) {
        return 1;
    }

    return aIndex - bIndex;
}