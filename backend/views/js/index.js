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
		/** @type {HTMLInputElement} */
		const $send_email = document.querySelector("#create-form #send_email");

		/** @type {BuyItem} */
		const newItem = {
			id: "",
			name: $name.value,
			current_quantity: Number($current_quantity.value),
			min_quantity: Number($min_quantity.value),
			send_email: $send_email.checked,
		};

		postBuyItem(newItem).then(_ => {
			fetchBuyItems().then(r => {
				loadItems("items-list", r);
			});
			fetchToBuyItems().then(r => {
				loadItems("to-buy-list", r);
			});
		});
	}, 
	true
);

/**
* @returns { Promise<BuyItem[]> }
* */
async function fetchToBuyItems() {
	const toBuyItems = await fetch("http://localhost:3000/to_buy_list");
	return await toBuyItems.json();
}

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

/** 
* @param {BuyItem[]} buyItems  
* @param {string} listId 
* */
function loadItems(listId, buyItems) {
	let $list = document.getElementById(listId);

	if ($list == null) {
		$list = document.createElement("ul");
		$list.id = listId;
		$root.appendChild($list);
	}
	$list.innerHTML = "";

	for (const { id, name, current_quantity, min_quantity, send_email} of buyItems) {
		const $item = document.createElement("li");
		const $div = document.createElement("div");
		const $checkbox = document.createElement("input");
		$checkbox.type = "checkbox";

		$item.innerText += `[${id}]${name}: ${current_quantity}/${min_quantity} [${send_email}]`;
		$div.style = 'display: flex;'
		$div.appendChild($checkbox);
		$div.appendChild($item);

		$list.appendChild($div);
	}
}

fetchBuyItems().then(
	/** @param {BuyItem[]} to_buy_list  */
	(totalItems) => {
		loadItems("items-list", totalItems);
	}
);

fetchToBuyItems().then(
	/** @param {BuyItem[]} to_buy_list  */
	(toBuyItems) => {
		loadItems("to-buy-list", toBuyItems);
	}
);
