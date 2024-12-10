
export function toggleSidebar(sidebar) {
    if (sidebar.classList.contains('open')) {
        sidebar.classList.remove('open');
    } else {
        sidebar.classList.add('open');
    }
}

export function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

export function timeAffi(time) {
    let now = new Date()
    let then = new Date(time)
    let min = ((now - then) / (1000 * 60)).toFixed()
    let hours = ((now - then) / (1000 * 3600)).toFixed()
    let days = ((now - then) / (1000 * 3600 * 24)).toFixed()
    if (min < 60) {
        return min + " minutes"
    } else if (hours < 24) {
        return hours + " hours"
    } else {
        return days + " day"
    }
}

export function FetchComments(commentForms) {
    commentForms.forEach(cForm => {
        cForm.addEventListener("submit", async (e) => {
            e.preventDefault()
            const path = e.currentTarget.getAttribute("action")
            const formData = new FormData(e.currentTarget)
            const TextA = e.currentTarget.querySelector("textarea")
            try {
                let res = await fetch(path, {
                    method: "post",
                    body: formData
                })
                
                if (!res.ok) {
                    switch (res.status) {
                        case 400:
                            throw new Error("Bad Request: Please check your input.");
                        case 401:
                            throw new Error("Forbidden: You need to log in to perform this action.");
                        case 404:
                            throw new Error("Not Found: The requested resource could not be found.");
                        case 500:
                            throw new Error("Server Error: Please try again later.");
                        default:
                            throw new Error(`Unexpected Error: ${res.status} ${res.statusText}`);
                    }
                }

                let json = await res.json()
                console.log("the json response here :", json);
                console.log("contentC :", TextA.value);
                const div = document.createElement("div")
                if (TextA.value != "") {
                    div.innerHTML = `<div style="color:green;"><p>Comment Added Succesfully!</p></div>`
                    cForm.appendChild(div)
                    setTimeout(() => {
                        div.remove()
                    }, 2000)
                }
                TextA.value = ""
            } catch (error) {
                alert(error)
                console.log(error);

                if (error.message == "Forbidden: You need to log in to perform this action.") {
                    window.location.href = '/login'
                }
            }
        })
    })

}

export function FetchReactions(Forms) {
    Forms.forEach(Form => {
        Form.addEventListener("submit", async (e) => {
            e.preventDefault();

            const button = e.currentTarget.querySelector('button');
            const choice = e.currentTarget.getAttribute('id')
            let TheOtherButton = ""
            if (choice == "like_post") {
                TheOtherButton = (e.currentTarget.parentElement.querySelector('#dislike_post').querySelector('button'))
    
            } else {
                TheOtherButton = (e.currentTarget.parentElement.querySelector('#like_post').querySelector('button'))
            }


            if (!button) {
                console.error("Button not found");
                return;
            }

            let path1 = e.currentTarget.getAttribute('action');

            let icon = "👍🏻"
            let icon2 = "👎🏻"
            // console.log("icon2", icon2)
            let previouslike = button.textContent
            console.log("previouslike outside ===>", previouslike)
            const postid = e.currentTarget.parentElement.querySelector('input').value;

            try {
                await fetch(path1, {
                    method: "POST",
                });

                const path = `/api/likes?postid=${postid}`;

                const response = await getdata(path);
                console.log(response);

                const likeCount = response.LikeCount;
                const DislikeCount = response.DislikeCount;

                if (choice == "like_post") {
                    button.textContent = icon + " " + likeCount;
                    TheOtherButton.textContent = icon2 + " " + DislikeCount
                    button.style.backgroundColor = "green"
                    TheOtherButton.style.backgroundColor = ""
                } else if (choice == "dislike_post") {
                    button.textContent = icon2 + " " + DislikeCount;
                    TheOtherButton.textContent = icon + " " + likeCount
                    button.style.backgroundColor = "red"
                    TheOtherButton.style.backgroundColor = ""
                }
            } catch (err) {
                console.error("Error processing like:", err);
            }
        });
    });
}

export function FetchCommentsReactions(Forms) {
    Forms.forEach((CurrForm) => {
        CurrForm.addEventListener('submit', async (event) => {
            event.preventDefault()
            const ActionPath = event.currentTarget.getAttribute('action')
            const Comment_id = event.currentTarget.parentElement.querySelector('input').value
            const choice = event.currentTarget.getAttribute('id')
            const button = event.currentTarget.querySelector('button')
            let previouslike = button.textContent
            console.log(" previouslikje ===>", previouslike);
            let TheOtherButton = ""
            if (choice == "Comment_Like") {
                TheOtherButton = event.currentTarget.parentElement.querySelector('#Comment_Dislike').querySelector('button')
            } else {
                TheOtherButton = event.currentTarget.parentElement.querySelector('#Comment_Like').querySelector('button')
            }

            let icon = button.textContent.split(' ')[0];
            let icon2 = TheOtherButton.innerHTML.split(' ')[0];
            console.log("iocnes :", icon, icon2);
            try {
                await fetch(ActionPath, {
                    method: "POST",
                })
                const ApiPath = `/api/likes?comment_id=${Comment_id}`;
                const response = await getdata(ApiPath)
                const likeCount = response.LikeCount
                const DislikeCount = response.DislikeCount;
                if (choice == "Comment_Like") {
                    button.textContent = icon + " " + likeCount;
                    TheOtherButton.textContent = icon2 + " " + DislikeCount
                    button.style.backgroundColor = "green"
                    TheOtherButton.style.backgroundColor = ""
                } else if (choice == "Comment_Dislike") {
                    button.textContent = icon + " " + DislikeCount;
                    TheOtherButton.textContent = icon2 + " " + likeCount
                    button.style.backgroundColor = "red"
                    TheOtherButton.style.backgroundColor = ""
                }
            } catch (error) {
                console.log("Error while reaction with commments", error)
            }
        })
    })
} 

async function getdata(path) {
    try {
        const res = await fetch(path);
        console.log(res);

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    throw new Error("Bad Request: Please check your input.");
                case 401:
                    throw new Error("Forbidden: You need to log in to perform this action.");
                case 404:
                    throw new Error("Not Found: The requested resource could not be found.");
                case 500:
                    throw new Error("Server Error: Please try again later.");
                default:
                    throw new Error(`Unexpected Error: ${res.status} ${res.statusText}`);
            }
        }

        let json = await res.json()
        return json
    } catch (error) {
        alert(error)
        if (error.message == "Forbidden: You need to log in to perform this action.") {
            window.location.href = '/login'
        }
    }
}

