// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

/**
 * Note : This test requires the demo plugin tar file under fixtures folder.
 * Download :
 * make dist latest master and copy to ./e2e/cypress/fixtures/com.mattermost.demo-plugin-0.9.0.tar.gz
 */

import {verifyEphemeralMessage} from 'mattermost-webapp/e2e/cypress/integration/integrations/builtin_commands/helper';

describe('Test App command autocomplete', () => {
    const pluginID = Cypress.config('pluginID');
    const pluginFile = Cypress.config('pluginFile');

    before(() => {
        cy.apiAdminLogin();
        cy.visit('/');

        // cy.apiRemovePluginById(pluginID);

        // cy.apiUploadPlugin(pluginFile);
        cy.apiEnablePluginById(pluginID);
    });

    after(() => {
        // cy.apiRemovePluginById(pluginID);
    });

    it('/http-hello message', () => {
        // # Type `/http-hello`
        cy.get('#post_textbox').clear().type('/http-hello');

        // * Verify autocomplete hint is correct
        checkSuggestions([
            {title: 'http-hello', hint: ''},
        ]);

        // # Type `/http-hello `
        cy.get('#post_textbox').type(' ');

        // * Check autocomplete options
        checkSuggestions([
            {title: 'message [--user] message', hint: 'send a message to a user'},
            {title: 'message-modal [--message] message', hint: 'send a message to a user'},
            {title: 'manage subscribe | unsubscribe', hint: 'manage channel subscriptions to greet new users'},
        ]);

        // # Choose "message"
        clickSuggestion(1);

        checkSuggestions([
            {title: '@ enter user ID or @user', hint: 'User to send the survey to'},
        ]);
        cy.get('#suggestionList > .slash-command:nth-child(1)').click();
        cy.wait(2000);

        cy.get('#post_textbox').type('{backspace}');
        cy.get('#post_textbox').type('sysadmin');
        cy.wait(2000);
        cy.get(`#suggestionList .mentions__name`).first().click();
        cy.get('#post_textbox').type(' ');

        checkSuggestions([
            {title: '--other Pick one', hint: 'Some values'},
            {title: '--message Anything you want to say', hint: 'Text to ask the user about'},
        ]);
        clickSuggestion(1);

        checkSuggestions([
            {title: 'option1', hint: 'Option 1'},
        ]);
        clickSuggestion(1);

        cy.get('#post_textbox').type(' ');

        checkSuggestions([
            {title: '--message Anything you want to say', hint: 'Text to ask the user about'},
        ]);
        clickSuggestion(1);

        cy.get('#post_textbox').type('themessage');
        cy.get('#post_textbox').type(' ');
        cy.get('#post_textbox').type('{enter}');

        verifyEphemeralMessage('Successfully sent survey');
    });
});

const checkSuggestions = (assertions) => {
    assertions.forEach((assertion, i) => {
        cy.get(`#suggestionList > .slash-command:nth-child(${i+1}) .slash-command__title`).should('contain.text', assertion.title);
        cy.get(`#suggestionList > .slash-command:nth-child(${i+1}) .slash-command__desc`).should('contain.text', assertion.hint);
    });
}

const clickSuggestion = (i) => {
    cy.get(`#suggestionList > .slash-command:nth-child(${i})`).click();
    cy.wait(2000);
    cy.get('#post_textbox').type(' ');
    cy.wait(2000);
}
