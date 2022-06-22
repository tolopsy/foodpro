import React from "react";

const Recipe = (props) => {
  return (
    <div class="recipe">
      <h4>{props.recipe.name}</h4>
      <ul>
        {
          props.recipe.ingredients && props.recipe.ingredients.map((ingredient, index) => <li>{ingredient}</li>)
        }
      </ul>
    </div>
  )
}

export default Recipe