const $root = document.createElement("div");
/**
* @typedef {Object} BuyItem
* @property {string} id
* @property {string} name
* @property {number} current_quantity
* @property {number} min_quantity
* @property {boolean} send_email
* */

$root.id = "root";

const $body = document.getElementsByTagName("body")[0];
$body.appendChild($root);

// Adding create buy item form
const $createForm = document.getElementById("create-form");
$createForm.addEventListener("submit",
	(e) => {
		e.preventDefault();
		const $name = document.querySelector("#create-form #name");
		const $current_quantity = document.querySelector("#create-form #current_quantity");
		const $min_quantity = document.querySelector("#create-form #min_quantity");
		const $send_email = document.querySelector("#create-form #send_email");

		/** @type {BuyItem} */
		const newItem = {
			id: "",
			name: $name.value,
			current_quantity: Number($current_quantity.value),
			min_quantity: Number($min_quantity.value),
			send_email: $send_email.value === "on",
		};

		postBuyItem(newItem).then(_ => {
			fetchBuyItems().then(r => {
				loadItems(r);
			});
		});
	}, 
	true
);

async function fetchBuyItems() {
	const buyItems = await fetch("http://localhost:3000/buy_items");
	return await buyItems.json();
}

/** @param {BuyItem} newBuyItem */
async function postBuyItem(newBuyItem) {
	let response = await fetch(
		"http://localhost:3000/buy_items", 
		{ 
			method: "POST", 
			body: JSON.stringify(newBuyItem) 
		}
	);

	return response;
}

/** @param {BuyItem[]} buyItems  */
function loadItems(buyItems) {
	let $list = document.getElementById("items-list");

	console.log($list);
	if ($list == null) {
		$list = document.createElement("ul");
		$list.id = "items-list";
		$root.appendChild($list);
	}
	$list.innerHTML = "";

	for (const { id, name, current_quantity, min_quantity, send_email} of buyItems) {
		const $item = document.createElement("li");
		$item.innerText = `[${id}]${name}: ${current_quantity}/${min_quantity} [${send_email}]`;
		$list.appendChild($item);
	}
}

fetchBuyItems().then(r => {
	console.log(r);
	loadItems(r);

	postBuyItem({ id: "3", name: "T4", current_quantity: 42, min_quantity: 1, send_email: true }).then(r => {
		console.log(r);

		if (r.status == 201) {
			fetchBuyItems().then(r => {
				loadItems(r);
			})
		}
	});
})


