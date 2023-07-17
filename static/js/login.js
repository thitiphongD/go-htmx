function openTab(event, tabName) {
    var i, tabContent;

    tabContent = document.getElementsByClassName('tab-content');
    for (i = 0; i < tabContent.length; i++) {
        tabContent[i].style.display = 'none';
    }

    var tabButtons = document.getElementsByClassName('tab-button');
    for (i = 0; i < tabButtons.length; i++) {
        tabButtons[i].classList.remove('active');
    }

    document.getElementById(tabName).style.display = 'block';

    event.currentTarget.classList.add('active');
}

document.getElementById('tab1').style.display = 'block';
document.getElementsByClassName('tab-button')[0].classList.add('active');
