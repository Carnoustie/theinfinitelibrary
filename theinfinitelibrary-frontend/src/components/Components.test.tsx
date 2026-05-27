import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { Home } from './Components';

test('renders home page with join prompt', () => {
    render(
        <BrowserRouter>
            <Home join_prompt="Join us!"/>
        </BrowserRouter>
    );
    expect(screen.getByText('Join us!')).toBeInTheDocument();
});