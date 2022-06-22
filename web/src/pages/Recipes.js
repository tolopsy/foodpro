import React, {useState, useEffect, useRef} from "react";
import Recipe from "../components/Recipe/Recipe";

const Recipes = () => {
  const [recipes, setRecipes] = useState([])
  const hasFetchedData = useRef(false)
  const getRecipes = () => {
    const endpoint = "http://localhost:8080/recipes"
    fetch(endpoint).then(
      response => response.json()
    ).then(
      data => setRecipes(data)
    ).catch(err => console.log(err))
  }

  useEffect(() => {
    // prevent double request in strict mode
    if (hasFetchedData.current === false) {
      getRecipes();
      hasFetchedData.current = true;
    } 
  }, [])

  return (
    <div>
      {recipes.map((recipe, index) => <Recipe recipe={recipe} />)}
    </div>
  )
}

export default Recipes